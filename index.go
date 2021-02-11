package kvas

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// indexPath computes filepath to a valueSet index
func (vs *ValueSet) indexPath() string {
	ip := filepath.Join(vs.baseDir, indexFilename+vs.indExt)
	return ip
}

// initIndex initializes index data structure
func (vs *ValueSet) initIndex() {
	if vs.index == nil {
		vs.index = make(map[string]IndexRecord, 0)
	}
}

// readIndex reads index of a valueSet
func (vs *ValueSet) readIndex() error {
	indexPath := vs.indexPath()

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
	}

	bytes, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &vs.index)
}

// writeIndex writes index of a valueSet
func (vs *ValueSet) writeIndex() error {
	bytes, err := json.Marshal(vs.index)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(vs.indexPath(), bytes, filePerm)
}

// setIndex updates index by key
func (vs *ValueSet) setIndex(key string, hash string) {
	if _, ok := vs.index[key]; !ok {
		vs.index[key] = IndexRecord{
			Created: time.Now().Unix(),
		}
	}

	if vs.index[key].Hash != hash {
		vs.index[key] = IndexRecord{
			Hash:     hash,
			Created:  vs.index[key].Created,
			Modified: time.Now().Unix(),
		}
	} else {
		log.Printf("ValueSet.setIndex: hash for item with key '%s' is the same", key)
	}
}
