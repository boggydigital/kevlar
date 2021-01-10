package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	jsonExt              = ".json"
	filePerm os.FileMode = 0644
	dirPerm  os.FileMode = 0755
)

type JsonValueSet struct {
	valueSet
}

func NewJsonClient(location string) (*JsonValueSet, error) {
	jvs := &JsonValueSet{valueSet{baseDir: location}}
	err := jvs.readIndex()
	return jvs, err
}

func (vs *JsonValueSet) indexPath() string {
	return filepath.Join(vs.baseDir, indexFilename+jsonExt)
}

func (vs *JsonValueSet) valuePath(key string) string {
	return filepath.Join(vs.baseDir, key+jsonExt)
}

func (vs *JsonValueSet) readIndex() error {
	indexPath := vs.indexPath()

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
	}

	bytes, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, &vs.index); err != nil {
		return err
	}

	vs.initIndex()

	return nil
}

func (vs *JsonValueSet) writeIndex() error {

	bytes, err := json.Marshal(vs.index)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(vs.indexPath(), bytes, filePerm)
}

func (vs *JsonValueSet) Get(key string) ([]byte, error) {
	if !vs.Contains(key) {
		return nil, nil
	}

	valuePath := vs.valuePath(key)
	if _, err := os.Stat(valuePath); os.IsNotExist(err) {
		return nil, nil
	}
	return ioutil.ReadFile(valuePath)
}

func (vs *JsonValueSet) Set(key string, value []byte) error {
	// check if value already exists and has the same hash
	hash, err := Sha256(value)
	if err != nil {
		return err
	}

	// initialize index if that's one first set operation
	//vs.initIndex()

	if hash == vs.index[key].Hash {
		return nil
	}

	// write value
	valuePath := vs.valuePath(key)
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
	vs.setIndex(key, hash)
	return vs.writeIndex()
}

func (vs *JsonValueSet) Remove(key string) error {
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
