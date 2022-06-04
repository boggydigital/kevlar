package kvas

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type localKeyValues struct {
	dir      string
	ext      string
	idx      index
	mtx      sync.Mutex
	connTime int64
}

const (
	JsonExt = ".json"
	GobExt  = ".gob"
)

const dirPerm os.FileMode = 0755

func ConnectLocal(dir string, ext string) (KeyValues, error) {

	switch ext {
	case JsonExt:
		fallthrough
	case GobExt:
		//do nothing
	default:
		return nil, fmt.Errorf("unknown extension %s", ext)
	}

	kv := &localKeyValues{
		dir: dir,
		ext: ext,
		idx: make(index, 0),
		mtx: sync.Mutex{},
	}
	err := kv.idx.read(kv.dir)

	kv.connTime = time.Now().Unix()

	return kv, err
}

// Has verifies if a value set contains provided key
func (lkv *localKeyValues) Has(key string) bool {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	_, ok := lkv.idx[key]
	return ok
}

func (lkv *localKeyValues) Get(key string) (io.ReadCloser, error) {

	if !lkv.Has(key) {
		return nil, nil
	}

	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	valAbsPath := lkv.valuePath(key)
	if _, err := os.Stat(valAbsPath); os.IsNotExist(err) {
		return nil, nil
	}
	return os.Open(valAbsPath)
}

func (lkv *localKeyValues) valuePath(key string) string {
	return filepath.Join(lkv.dir, key+lkv.ext)
}

// Set stores a bytes slice value by a provided key
func (lkv *localKeyValues) Set(key string, reader io.Reader) error {

	var buf bytes.Buffer
	tr := io.TeeReader(reader, &buf)

	// check if value already exists and has the same hash
	hash, err := Sha256(tr)
	if err != nil {
		return err
	}

	lkv.mtx.Lock()

	if hash == lkv.idx[key].Hash {
		lkv.mtx.Unlock()
		return nil
	}

	lkv.mtx.Unlock()

	// write value
	valuePath := lkv.valuePath(key)

	if _, err := os.Stat(lkv.dir); os.IsNotExist(err) {
		if err := os.MkdirAll(lkv.dir, dirPerm); err != nil {
			return err
		}
	}
	file, err := os.Create(valuePath)
	defer file.Close()
	if err != nil {
		return err
	}

	if _, err = io.Copy(file, &buf); err != nil {
		return err
	}

	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	// update index
	lkv.idx.upd(key, hash)
	return lkv.idx.write(lkv.dir)
}

// Cut deletes value from localKeyValues by a provided key
func (lkv *localKeyValues) Cut(key string) (bool, error) {

	if !lkv.Has(key) {
		return false, nil
	}

	// delete value
	valuePath := lkv.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		return false, fmt.Errorf("index contains key %s, file not found", key)
	}

	if err := os.Remove(valuePath); err != nil {
		return false, err
	}

	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	// update index
	delete(lkv.idx, key)

	return true, lkv.idx.write(lkv.dir)
}

func (lkv *localKeyValues) Keys() []string {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	return lkv.idx.Keys(lkv.mtx)
}

// CreatedAfter returns keys of values created on or after provided timestamp
func (lkv *localKeyValues) CreatedAfter(timestamp int64) []string {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	return lkv.idx.CreatedAfter(timestamp, lkv.mtx)
}

// ModifiedAfter returns keys of values modified on or after provided timestamp
// that were created earlier
func (lkv *localKeyValues) ModifiedAfter(timestamp int64, strictlyModified bool) []string {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	return lkv.idx.ModifiedAfter(timestamp, strictlyModified, lkv.mtx)
}

func (lkv *localKeyValues) IsModifiedAfter(key string, timestamp int64) bool {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	return lkv.idx.IsModifiedAfter(key, timestamp, lkv.mtx)
}

func (lkv *localKeyValues) IndexCurrentModTime() (int64, error) {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	indexPath := indexPath(lkv.dir)
	if stat, err := os.Stat(indexPath); os.IsNotExist(err) {
		return -1, nil
	} else if err != nil {
		return -1, err
	} else {
		return stat.ModTime().Unix(), nil
	}
}

func (lkv *localKeyValues) CurrentModTime(key string) (int64, error) {
	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	valuePath := lkv.valuePath(key)
	if stat, err := os.Stat(valuePath); os.IsNotExist(err) {
		return -1, nil
	} else if err != nil {
		return -1, err
	} else {
		return stat.ModTime().Unix(), nil
	}
}

func (lkv *localKeyValues) IndexRefresh() error {
	indexModTime, err := lkv.IndexCurrentModTime()
	if err != nil {
		return err
	}

	lkv.mtx.Lock()
	defer lkv.mtx.Unlock()

	if lkv.connTime < indexModTime {
		if err := lkv.idx.read(lkv.dir); err != nil {
			return err
		}
	}

	return nil
}
