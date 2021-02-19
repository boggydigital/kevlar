package kvas

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"time"
)

// indexPath computes filepath to a valueSet index
func (vs *ValueSet) indexPath() string {
	ip := filepath.Join(vs.baseDir, indexFilename+gobExt)
	return ip
}

//// initIndex initializes index data structure
//func (vs *ValueSet) initIndex() {
//if vs.index == nil {
//	vs.index = make(map[string]IndexRecord, 0)
//}
//}

// readIndex reads index of a valueSet
func (vs *ValueSet) readIndex() error {
	indexPath := vs.indexPath()

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
	}

	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	return gob.NewDecoder(indexFile).Decode(&vs.index)
}

// writeIndex writes index of a valueSet
func (vs *ValueSet) writeIndex() error {

	indexFile, err := os.Create(vs.indexPath())
	if err != nil {
		return err
	}
	defer indexFile.Close()

	return gob.NewEncoder(indexFile).Encode(vs.index)
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
