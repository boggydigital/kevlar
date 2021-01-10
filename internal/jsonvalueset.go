package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// JsonValueSet is a valueSet that stores values in .json files
type JsonValueSet struct {
	valueSet
}

// NewJsonClient creates a JsonValueSet client at the provided
// location, where the index and the values would be stored
func NewJsonClient(location string) (*JsonValueSet, error) {
	jvs := &JsonValueSet{valueSet{baseDir: location}}
	err := jvs.readIndex()
	return jvs, err
}

// getExt returns extension of the files in valueSet
func (*JsonValueSet) getExt() string {
	return jsonExt
}

// indexPath computes filepath to a valueSet index
func (jvs *JsonValueSet) indexPath() string {
	return filepath.Join(jvs.baseDir, indexFilename+jvs.getExt())
}

// valuePath computes filepath to a value by key
func (jvs *JsonValueSet) valuePath(key string) string {
	return filepath.Join(jvs.baseDir, key+jvs.getExt())
}

// readIndex reads index of a valueSet
func (jvs *JsonValueSet) readIndex() error {
	indexPath := jvs.indexPath()

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
	}

	bytes, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &jvs.index)
}

// writeIndex writes index of a valueSet
func (jvs *JsonValueSet) writeIndex() error {
	bytes, err := json.Marshal(jvs.index)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(jvs.indexPath(), bytes, filePerm)
}

// Get returns a bytes slice value by a provided key
func (jvs *JsonValueSet) Get(key string) ([]byte, error) {
	if !jvs.Contains(key) {
		return nil, nil
	}

	valuePath := jvs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		return nil, nil
	}
	return ioutil.ReadFile(valuePath)
}

// Set stores a bytes slice value by a provided key
func (jvs *JsonValueSet) Set(key string, value []byte) error {
	// check if value already exists and has the same hash
	hash, err := Sha256(value)
	if err != nil {
		return err
	}

	if hash == jvs.index[key].Hash {
		return nil
	}

	// write value
	valuePath := jvs.valuePath(key)
	dirPath := filepath.Dir(valuePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, dirPerm)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(valuePath, value, filePerm)
	if err != nil {
		return err
	}

	// update index
	jvs.initIndex()
	jvs.setIndex(key, hash)
	return jvs.writeIndex()
}

// Remove deletes value from a valueSet by a provided key
func (jvs *JsonValueSet) Remove(key string) error {
	if !jvs.Contains(key) {
		return nil
	}

	// delete value
	valuePath := jvs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(valuePath); err != nil {
		return err
	}

	// update index
	jvs.initIndex()
	delete(jvs.index, key)
	return jvs.writeIndex()
}
