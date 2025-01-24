package kevlar_legacy

import (
	"bytes"
	"encoding/gob"
	"github.com/boggydigital/busan"
	"golang.org/x/exp/maps"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	kevlarDirname      = "_kevlar"
	logRecordsFilename = "_log.gob"
	hashExt            = ".sha256"
)

type keyValues struct {
	dir  string
	ext  string
	lmt  int64
	log  logRecords
	keys map[string]any
	mtx  *sync.Mutex
}

// NewKeyValues connects a new local key value storage at the specified directory
// and will use specified extension for the value files
func NewKeyValues(dir, ext string) (KeyValues, error) {

	// make sure dir we're connecting to exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	kv := &keyValues{
		dir: dir,
		ext: ext,
		mtx: new(sync.Mutex),
	}

	_, kv.lmt = kv.IsCurrent()

	if err := kv.refreshLogRecords(); os.IsNotExist(err) {
		// do nothing
	} else if err != nil {
		return nil, err
	}

	return kv, nil
}

func (kv *keyValues) IsCurrent() (bool, int64) {
	var lmt int64 = -1
	if fi, err := os.Stat(kv.absLogRecordsFilename()); err == nil {
		lmt = fi.ModTime().Unix()
	}
	return lmt == kv.lmt, lmt
}

func (kv *keyValues) refreshLogRecords() error {
	if ok, lmt := kv.IsCurrent(); ok {
		if kv.log != nil {
			return nil
		}
	} else {
		kv.mtx.Lock()
		kv.lmt = lmt
		kv.mtx.Unlock()
	}

	absLogFilename := kv.absLogRecordsFilename()
	if _, err := os.Stat(absLogFilename); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	logFile, err := os.Open(absLogFilename)
	if err != nil {
		return err
	}
	defer logFile.Close()

	kv.mtx.Lock()
	defer kv.mtx.Unlock()

	if err := gob.NewDecoder(logFile).Decode(&kv.log); err == io.EOF {
		// do nothing - empty log will be initialized later
	} else if err != nil {
		return err
	}

	return nil
}

func (kv *keyValues) refreshKeys() error {

	if err := kv.refreshLogRecords(); err != nil {
		return err
	}

	uks := make(map[string]any)
	for _, lr := range kv.log {
		switch lr.Mt {
		case create:
			fallthrough
		case update:
			uks[lr.Id] = nil
		case cut:
			delete(uks, lr.Id)
		default:
			panic("unknown log record mutation type")
		}
	}

	kv.mtx.Lock()
	defer kv.mtx.Unlock()

	kv.keys = uks

	return nil
}

func (kv *keyValues) Keys() ([]string, error) {
	if err := kv.refreshKeys(); err != nil {
		return nil, err
	}

	kv.mtx.Lock()
	defer kv.mtx.Unlock()

	return maps.Keys(kv.keys), nil
}

func (kv *keyValues) Has(key string) (bool, error) {
	if err := kv.refreshKeys(); err != nil {
		return false, err
	}

	kv.mtx.Lock()
	defer kv.mtx.Unlock()

	_, ok := kv.keys[key]

	return ok, nil
}

func (kv *keyValues) absLogRecordsFilename() string {
	return filepath.Join(kv.dir, kevlarDirname, logRecordsFilename)
}

func (kv *keyValues) absValueFilename(key string) string {
	return filepath.Join(kv.dir, busan.Sanitize(key)+kv.ext)
}

func (kv *keyValues) absHashFilename(key string) string {
	return filepath.Join(kv.dir, kevlarDirname, busan.Sanitize(key)+hashExt)
}

func (kv *keyValues) Get(key string) (io.ReadCloser, error) {
	return os.Open(kv.absValueFilename(key))
}

func (kv *keyValues) currentHash(key string) (string, error) {
	if ok, err := kv.Has(key); err == nil {
		if !ok {
			return "", nil
		}
	} else {
		return "", err
	}

	absHashFilename := kv.absHashFilename(key)
	if _, err := os.Stat(absHashFilename); err != nil {
		return "", nil
	}
	hashFile, err := os.Open(absHashFilename)
	if err != nil {
		return "", err
	}
	defer hashFile.Close()

	sb := new(strings.Builder)

	if _, err := io.Copy(sb, hashFile); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (kv *keyValues) createLogRecords() error {
	absLogRecordsFilename := kv.absLogRecordsFilename()
	dir, _ := filepath.Split(absLogRecordsFilename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	logFile, err := os.Create(absLogRecordsFilename)
	if err != nil {
		return err
	}
	defer logFile.Close()

	if err := lockFd(logFile.Fd()); err != nil {
		return err
	}

	if err := gob.NewEncoder(logFile).Encode(kv.log); err != nil {
		return err
	}

	return unlockFd(logFile.Fd())
}

func (kv *keyValues) appendLogRecord(rec *logRecord) error {
	if err := kv.refreshLogRecords(); err != nil {
		return err
	}

	kv.mtx.Lock()
	defer kv.mtx.Unlock()

	kv.log = append(kv.log, rec)

	return kv.createLogRecords()

}

func (kv *keyValues) createLogRecord(key string) error {
	// adding the key right away to respond to Has queries before log update
	kv.mtx.Lock()
	kv.keys[key] = nil
	kv.mtx.Unlock()

	rec := &logRecord{
		Ts: time.Now().Unix(),
		Mt: create,
		Id: key,
	}

	return kv.appendLogRecord(rec)
}

func (kv *keyValues) updateLogRecord(key string) error {
	kv.mtx.Lock()
	updated := false
	for _, rec := range kv.log {
		if rec.Id == key && rec.Mt == update {
			rec.Ts = time.Now().Unix()
			updated = true
			break
		}
	}
	kv.mtx.Unlock()

	if updated {
		return kv.createLogRecords()
	} else {
		rec := &logRecord{
			Ts: time.Now().Unix(),
			Mt: update,
			Id: key,
		}
		return kv.appendLogRecord(rec)
	}
}

func (kv *keyValues) createOrUpdateLogRecord(key string) error {
	if ok, err := kv.Has(key); err == nil {
		if ok {
			return kv.updateLogRecord(key)
		} else {
			return kv.createLogRecord(key)
		}
	} else {
		return err
	}
}

func (kv *keyValues) cutLogRecord(key string) error {
	rec := &logRecord{
		Ts: time.Now().Unix(),
		Mt: cut,
		Id: key,
	}

	kv.mtx.Lock()
	delete(kv.keys, key)
	kv.mtx.Unlock()

	return kv.appendLogRecord(rec)
}

func (kv *keyValues) createHashFile(key, hash string) error {
	absHashFilename := kv.absHashFilename(key)
	dir, _ := filepath.Split(absHashFilename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	hashFile, err := os.Create(absHashFilename)
	if err != nil {
		return err
	}
	defer hashFile.Close()

	if _, err := io.Copy(hashFile, strings.NewReader(hash)); err != nil {
		return err
	}

	return nil
}

// Set writes the value to storage if the value has changed since the
// last time it was written. This is validated with a SHA-256 hash that
// is stored alongside the value in storage
func (kv *keyValues) Set(key string, reader io.Reader) error {

	var buf bytes.Buffer
	tr := io.TeeReader(reader, &buf)

	// check if value already exists and has the same hash
	hash, err := Sha256(tr)
	if err != nil {
		return err
	}

	currentHash, err := kv.currentHash(key)
	if err != nil {
		return err
	}

	// the latest value is already set
	if hash == currentHash {
		return nil
	}

	if err := kv.createHashFile(key, hash); err != nil {
		return err
	}

	// write value
	file, err := os.Create(kv.absValueFilename(key))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, &buf); err != nil {
		return err
	}

	return kv.createOrUpdateLogRecord(key)
}

// Cut removes the value from storage in the following sequence of events:
// - cut operation log value is added
// - stored hash value is removed
// - stored value is removed
func (kv *keyValues) Cut(key string) (bool, error) {
	if ok, err := kv.Has(key); err == nil {
		if !ok {
			return false, nil
		}
	} else {
		return false, err
	}

	absHashFilename := kv.absHashFilename(key)
	if _, err := os.Stat(absHashFilename); err == nil {
		if err := os.Remove(absHashFilename); err != nil {
			return false, err
		}
	}

	absValueFilename := kv.absValueFilename(key)
	if _, err := os.Stat(absValueFilename); err == nil {
		if err := os.Remove(absValueFilename); err != nil {
			return false, err
		}
	}

	if err := kv.cutLogRecord(key); err != nil {
		return false, err
	}

	return true, nil
}

func (kv *keyValues) filterLog(m func(*logRecord) bool) ([]string, error) {
	if err := kv.refreshLogRecords(); err != nil {
		return nil, err
	}
	matches := make(map[string]any)
	for _, lr := range kv.log {
		if m(lr) {
			matches[lr.Id] = nil
		}
		if lr.Mt == cut {
			delete(matches, lr.Id)
		}
	}
	return maps.Keys(matches), nil
}

func (kv *keyValues) CreatedAfter(ts int64) ([]string, error) {
	return kv.filterLog(func(r *logRecord) bool {
		return r.Mt == create && r.Ts >= ts
	})
}

func (kv *keyValues) UpdatedAfter(ts int64) ([]string, error) {
	return kv.filterLog(func(r *logRecord) bool {
		return r.Mt == update && r.Ts >= ts
	})
}

func (kv *keyValues) CreatedOrUpdatedAfter(ts int64) ([]string, error) {
	return kv.filterLog(func(r *logRecord) bool {
		createdAfter := r.Mt == create && r.Ts >= ts
		updatedAfter := r.Mt == update && r.Ts >= ts
		return createdAfter || updatedAfter
	})
}

func (kv *keyValues) IsUpdatedAfter(key string, ts int64) (bool, error) {
	filtered, err := kv.filterLog(func(r *logRecord) bool {
		if r.Id != key {
			return false
		}
		return r.Mt == update && r.Ts >= ts
	})
	if err != nil {
		return false, err
	}
	return len(filtered) > 0, nil
}

func (kv *keyValues) ModTime(key string) (int64, error) {
	if fi, err := os.Stat(kv.absValueFilename(key)); err == nil {
		return fi.ModTime().Unix(), nil
	} else if os.IsNotExist(err) {
		// key could have been deleted - check the log
		for _, lr := range kv.log {
			if lr.Id == key && lr.Mt == cut {
				return lr.Ts, nil
			}
		}
		return -1, nil
	} else {
		return -1, err
	}
}
