package kevlar

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/boggydigital/busan"
	"io"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

const (
	logRecordsFilename = "_log.gob"
	newExt             = ".new"
)

const (
	JsonExt = ".json"
	GobExt  = ".gob"
	HtmlExt = ".html"
	XmlExt  = ".xml"
)

const UnknownModTime = -1

type mutationType int

const (
	create mutationType = iota
	update
	cut
)

type logRecord struct {
	Id   string
	Ts   int64
	Mt   mutationType
	Hash []byte
}

type logRecords []*logRecord

type keyValues struct {
	dir string
	ext string
	log logRecords
	mtx *sync.Mutex
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

	if err := kv.loadLogRecords(); os.IsNotExist(err) {
		// do nothing, connecting to an empty key value store
	} else if err != nil {
		return nil, err
	}

	return kv, nil
}

func timeNow() int64 {
	return time.Now().UTC().Unix()
}

func createWriteOnlyFile(path string) (*os.File, error) {
	// not using O_EXCL intentionally here (meaning new file will be created even if the old exists)
	// existing file presence would indicate an incomplete write (crash) during previous operation
	// which among other things would mean:
	// - existing log is in good condition and is only missing that last attempted operation
	// - it's unclear what state existing file is, so it's not worth trying to salvage it
	// - instead we're just ignoring it to avoid blocking (hopefully) good operations
	dir, _ := filepath.Split(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
}

func sha256Bytes(reader io.Reader) ([]byte, error) {
	h := sha256.New()
	_, err := io.Copy(h, reader)
	return h.Sum(nil), err
}

func (kv *keyValues) loadLogRecords() error {

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
	if err := gob.NewDecoder(logFile).Decode(&kv.log); err == io.EOF {
		// do nothing - empty log will be initialized later
	} else if err != nil {
		return err
	}
	kv.mtx.Unlock()

	return nil
}

func (kv *keyValues) absLogRecordsFilename() string {
	return filepath.Join(kv.dir, logRecordsFilename)
}

func (kv *keyValues) absValueFilename(key string) string {
	return filepath.Join(kv.dir, busan.Sanitize(key)+kv.ext)
}

func (kv *keyValues) writeAtomically(path string, r io.Reader) error {

	newPath := path + newExt

	newFile, err := createWriteOnlyFile(newPath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	if _, err := io.Copy(newFile, r); err != nil {
		return err
	}

	if err := newFile.Sync(); err != nil {
		return err
	}

	return os.Rename(newPath, path)
}

// createLogRecord appends a new create log record
func (kv *keyValues) createLogRecord(key string, hash []byte) error {
	rec := &logRecord{
		Id:   key,
		Ts:   timeNow(),
		Mt:   create,
		Hash: hash,
	}

	return kv.writeLogRecord(rec)
}

// updateLogRecord removes all existing log records of type update and
// appends a new update log record
func (kv *keyValues) updateLogRecord(key string, hash []byte) error {
	kv.mtx.Lock()
	compactedLogRecords := make(logRecords, 0, len(kv.log))
	for _, lr := range kv.log {
		if lr.Id == key && lr.Mt == update {
			continue
		}
		compactedLogRecords = append(compactedLogRecords, lr)
	}

	kv.log = compactedLogRecords
	kv.mtx.Unlock()

	updLr := &logRecord{
		Id:   key,
		Ts:   timeNow(),
		Mt:   update,
		Hash: hash,
	}

	return kv.writeLogRecord(updLr)
}

// cutLogRecord removes all existing log records (any type) for this key and
// appends a new cut log record
func (kv *keyValues) cutLogRecord(key string) error {
	kv.mtx.Lock()
	compactedLogRecords := make(logRecords, 0, len(kv.log))
	for _, lr := range kv.log {
		if lr.Id == key {
			continue
		}
		compactedLogRecords = append(compactedLogRecords, lr)
	}

	kv.log = compactedLogRecords
	kv.mtx.Unlock()

	rec := &logRecord{
		Id: key,
		Ts: timeNow(),
		Mt: cut,
	}

	return kv.writeLogRecord(rec)
}

func (kv *keyValues) writeLogRecord(rec *logRecord) error {

	kv.mtx.Lock()
	if rec != nil {
		kv.log = append(kv.log, rec)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv.log); err != nil {
		return err
	}

	if err := kv.writeAtomically(kv.absLogRecordsFilename(), &buf); err != nil {
		return err
	}
	kv.mtx.Unlock()

	return nil
}

func (kv *keyValues) currentHash(key string) []byte {
	for ii := len(kv.log) - 1; ii >= 0; ii-- {
		if lr := kv.log[ii]; lr.Id == key {
			// that should work even if the log record is cut type
			return lr.Hash
		}
	}

	return nil
}

func (kv *keyValues) keys() map[string]any {
	keys := make(map[string]any)

	for _, lr := range kv.log {
		if lr.Mt == cut {
			delete(keys, lr.Id)
			continue
		}
		keys[lr.Id] = nil
	}

	return keys
}

func (kv *keyValues) filterLog(m func(*logRecord) bool) iter.Seq[string] {
	matches := make(map[string]any)
	for _, lr := range kv.log {
		if m(lr) {
			matches[lr.Id] = nil
		}
		if lr.Mt == cut {
			delete(matches, lr.Id)
		}
	}
	return maps.Keys(matches)
}

func (kv *keyValues) Len() int {
	return len(kv.keys())
}

func (kv *keyValues) Keys() iter.Seq[string] {
	return maps.Keys(kv.keys())
}

func (kv *keyValues) Has(key string) bool {
	for _, lr := range kv.log {
		if lr.Id == key {
			return true
		}
	}
	return false
}

func (kv *keyValues) Get(key string) (io.ReadCloser, error) {
	return os.Open(kv.absValueFilename(key))
}

// Set writes the value to storage if the value has changed since the
// last time it was written. This is validated with a SHA-256 hash that
// is stored in log
func (kv *keyValues) Set(key string, reader io.Reader) error {

	var buf bytes.Buffer
	tr := io.TeeReader(reader, &buf)

	// check if value already exists and has the same hash
	hash, err := sha256Bytes(tr)
	if err != nil {
		return err
	}

	currentHash := kv.currentHash(key)

	if slices.Equal(hash, currentHash) {
		return nil
	}

	kv.mtx.Lock()
	if err := kv.writeAtomically(kv.absValueFilename(key), &buf); err != nil {
		return err
	}
	kv.mtx.Unlock()

	if kv.Has(key) {
		return kv.updateLogRecord(key, hash)
	} else {
		return kv.createLogRecord(key, hash)
	}
}

// Cut removes the value from storage in the following sequence of events:
// - cut operation log value is added
// - stored value is removed
func (kv *keyValues) Cut(key string) (bool, error) {
	if !kv.Has(key) {
		return false, nil
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

func (kv *keyValues) CreatedAfter(ts int64) iter.Seq[string] {
	return kv.filterLog(func(r *logRecord) bool {
		return r.Mt == create && r.Ts >= ts
	})
}

func (kv *keyValues) UpdatedAfter(ts int64) iter.Seq[string] {
	return kv.filterLog(func(r *logRecord) bool {
		return r.Mt == update && r.Ts >= ts
	})
}

func (kv *keyValues) CreatedOrUpdatedAfter(ts int64) iter.Seq[string] {
	return kv.filterLog(func(r *logRecord) bool {
		createdAfter := r.Mt == create && r.Ts >= ts
		updatedAfter := r.Mt == update && r.Ts >= ts
		return createdAfter || updatedAfter
	})
}

func (kv *keyValues) IsUpdatedAfter(key string, ts int64) bool {
	for ii := len(kv.log) - 1; ii >= 0; ii-- {
		if lr := kv.log[ii]; lr.Id == key {
			switch lr.Mt {
			case update:
				return lr.Ts >= ts
			default:
				return false
			}
		}
	}
	return false
}

func (kv *keyValues) ModTime() int64 {
	if len(kv.log) > 0 {
		return kv.log[len(kv.log)-1].Ts
	}
	return UnknownModTime
}

func (kv *keyValues) ValueModTime(key string) int64 {
	for ii := len(kv.log) - 1; ii >= 0; ii-- {
		if lr := kv.log[ii]; lr.Id == key {
			return lr.Ts
		}
	}
	return UnknownModTime
}
