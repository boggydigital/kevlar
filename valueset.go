package kvas

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type ValueSet struct {
	baseDir string
	valExt  string
	index   map[string]IndexRecord
}

const (
	jsonExt = ".json"
	gobExt  = ".gob"
)

// ValueSetClient defines functions of a key value store client
type ValueSetClient interface {
	// index
	indexPath() string
	//initIndex()
	readIndex() error
	writeIndex() error
	removeIndex(string)
	setIndex(string, string)
	// values
	valuePath(string) string
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Remove(string) error
	Contains(string) bool
	// enumerations
	reduce(func(IndexRecord) bool) []string
	All() []string
	CreatedAfter(int64) []string
	ModifiedAfter(int64) []string
}

// NewLocal creates a ValueSet client at the provided
// location, where the index and the values would be stored
func newLocal(dst string, valExt string) (*ValueSet, error) {
	vs := &ValueSet{
		baseDir: dst,
		valExt:  valExt,
		index:   make(map[string]IndexRecord, 0),
	}
	err := vs.readIndex()
	return vs, err
}

func NewJsonLocal(dst string) (*ValueSet, error) {
	return newLocal(dst, jsonExt)
}

func NewGobLocal(dst string) (*ValueSet, error) {
	return newLocal(dst, gobExt)
}

// valuePath computes filepath to a value by key
func (vs *ValueSet) valuePath(key string) string {
	vp := filepath.Join(vs.baseDir, key+vs.valExt)
	return vp
}

// Get returns a bytes slice value by a provided key
func (vs *ValueSet) Get(key string) (io.ReadCloser, error) {
	if !vs.Contains(key) {
		return nil, nil
	}

	valuePath := vs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		return nil, nil
	}
	return os.Open(valuePath)
}

// Set stores a bytes slice value by a provided key
func (vs *ValueSet) Set(key string, reader io.Reader) error {

	var buf bytes.Buffer
	tr := io.TeeReader(reader, &buf)

	// check if value already exists and has the same hash
	hash, err := Sha256(tr)
	if err != nil {
		return err
	}

	if hash == vs.index[key].Hash {
		return nil
	}

	// write value
	valuePath := vs.valuePath(key)
	dirPath := filepath.Dir(valuePath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, dirPerm); err != nil {
			return err
		}
	}
	file, err := os.Create(valuePath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(file, &buf); err != nil {
		return err
	}

	// update index
	vs.setIndex(key, hash)
	return vs.writeIndex()
}

// Remove deletes value from a valueSet by a provided key
func (vs *ValueSet) Remove(key string) error {
	if !vs.Contains(key) {
		return nil
	}

	// delete value
	valuePath := vs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(valuePath); err != nil {
		return err
	}

	// update index
	delete(vs.index, key)
	return vs.writeIndex()
}

// Contains verifies if a value set contains provided key
func (vs *ValueSet) Contains(key string) bool {
	if _, ok := vs.index[key]; ok {
		return ok
	}
	return false
}

//
//func (vs *ValueSet) Title(key string) string {
//	if !vs.Contains(key) {
//		return ""
//	}
//
//	return vs.index[key].Title
//}
//
//func (vs *ValueSet) SetTitle(key, title string) error {
//	if !vs.Contains(key) {
//		return nil
//	}
//
//	vs.setTitle(key, title)
//	return vs.writeIndex()
//}
