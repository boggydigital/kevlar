package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type IndexRecord struct {
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type ValueSet struct {
	baseDir string
	ext     string
	debug   bool
	index   map[string]IndexRecord
}

// ValueSetClient defines functions of a key value store client
type ValueSetClient interface {
	// index
	indexPath() string
	initIndex()
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

// NewJsonClient creates a ValueSet client at the provided
// location, where the index and the values would be stored
func NewClient(location string, ext string, debug bool) (*ValueSet, error) {
	vs := &ValueSet{baseDir: location, ext: ext, debug: debug}
	err := vs.readIndex()
	return vs, err
}

// valuePath computes filepath to a value by key
func (vs *ValueSet) valuePath(key string) string {
	vp := filepath.Join(vs.baseDir, key+vs.ext)
	if vs.debug {
		log.Printf("ValueSet.valuePath: '%s' => %s\n", key, vp)
	}
	return vp
}

// Get returns a bytes slice value by a provided key
func (vs *ValueSet) Get(key string) ([]byte, error) {
	if !vs.Contains(key) {
		if vs.debug {
			log.Printf("ValueSet.Get: no item with key '%s'\n", key)
		}
		return nil, nil
	}

	valuePath := vs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		if vs.debug {
			log.Printf("ValueSet.Get: index contains key %s, but file doesn't exist\n", key)
		}
		return nil, nil
	}
	return ioutil.ReadFile(valuePath)
}

// Set stores a bytes slice value by a provided key
func (vs *ValueSet) Set(key string, value []byte) error {
	// check if value already exists and has the same hash
	hash, err := Sha256(value)
	if err != nil {
		if vs.debug {
			log.Printf("ValueSet.Set: error getting hash for value with a key '%s'\n", key)
		}
		return err
	}

	if hash == vs.index[key].Hash {
		if vs.debug {
			log.Printf("ValueSet.Set: hash is the same for the item with the key '%s'\n", key)
		}
		return nil
	}

	// write value
	valuePath := vs.valuePath(key)
	dirPath := filepath.Dir(valuePath)
	if vs.debug {
		log.Println("ValueSet.Set: target directory to write value to is", dirPath)
	}
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, dirPerm)
		if err != nil {
			if vs.debug {
				log.Println("ValueSet.Set: error making directory:", err.Error())
			}
			return err
		}
	}
	err = ioutil.WriteFile(valuePath, value, filePerm)
	if err != nil {
		if vs.debug {
			log.Println("ValueSet.Set: error writing file:", err.Error())
		}
		return err
	}

	// update index
	vs.initIndex()
	vs.setIndex(key, hash)
	return vs.writeIndex()
}

// Remove deletes value from a valueSet by a provided key
func (vs *ValueSet) Remove(key string) error {
	if !vs.Contains(key) {
		if vs.debug {
			log.Println("ValueSet.Remove: there is no value with the key", key)
		}
		return nil
	}

	// delete value
	valuePath := vs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		if vs.debug {
			log.Printf("ValueSet.Remove: index contains key '%s', but the file doesn't exist\n", key)
		}
		return nil
	}

	if err := os.Remove(valuePath); err != nil {
		if vs.debug {
			log.Println("ValueSet.Remove: error removing the file:", err.Error())
		}
		return err
	}

	// update index
	vs.initIndex()
	delete(vs.index, key)
	return vs.writeIndex()
}

// Contains verifies if a value set contains provided key
func (vs *ValueSet) Contains(key string) bool {
	_, ok := vs.index[key]
	if vs.debug {
		doesExist := "does"
		if !ok {
			doesExist = "doesn't"
		}
		log.Printf("ValueSet.Contains: item with the key '%s' %s exist\n", key, doesExist)
	}
	return ok
}
