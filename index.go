package kvas

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"time"
)

const indexFilename = "_index" + GobExt

type record struct {
	//Title    string `json:"title"`
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type index map[string]record

func (idx index) read(dir string) error {

	indexPath := filepath.Join(dir, indexFilename)

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
		//return fmt.Errorf("index %s doesn't exist", indexPath)
	}

	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	return gob.NewDecoder(indexFile).Decode(&idx)
}

func (idx index) write(dir string) error {

	indexPath := filepath.Join(dir, indexFilename)

	indexFile, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	return gob.NewEncoder(indexFile).Encode(idx)
}

func (idx index) upd(key string, hash string) {

	if _, ok := idx[key]; !ok {
		idx[key] = record{
			Created: time.Now().Unix(),
		}
	}

	if idx[key].Hash == hash {
		return
	}

	idx[key] = record{
		Hash:     hash,
		Created:  idx[key].Created,
		Modified: time.Now().Unix(),
	}
}
