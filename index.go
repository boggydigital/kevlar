package kvas

import (
	"encoding/gob"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type IndexRecord struct {
	Title    string `json:"title"`
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

// indexPath computes filepath to a valueSet index
func (vs *ValueSet) indexPath() string {
	ip := filepath.Join(vs.baseDir, indexFilename+gobExt)
	return ip
}

// readIndex reads index of a valueSet
func (vs *ValueSet) readIndex() error {
	indexPath := vs.indexPath()

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
		//return fmt.Errorf("index %s doesn't exist", indexPath)
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

	mu := sync.Mutex{}
	defer mu.Unlock()
	mu.Lock()

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
