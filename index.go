package kvas

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const indexFilename = "_index" + GobExt

type record struct {
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type index map[string]record

func indexPath(dir string) string {
	return filepath.Join(dir, indexFilename)
}

func (idx index) read(dir string) error {

	indexPath := indexPath(dir)

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return nil
		//return fmt.Errorf("index %s doesn't exist", indexPath)
	}

	indexFile, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	if err := gob.NewDecoder(indexFile).Decode(&idx); err != nil {
		// attempt to gracefully recover index from the filesystem state
		//idx = make(index)

		di, err := os.Open(dir)
		if err != nil {
			return err
		}

		files, err := di.Readdir(-1)
		if err != nil {
			return err
		}

		idx = make(index)

		for _, fi := range files {
			fn := fi.Name()
			if fn == indexFilename {
				continue
			}

			fn = strings.TrimSuffix(fn, filepath.Ext(fn))

			idx[fn] = record{
				Created:  fi.ModTime().Unix(),
				Modified: fi.ModTime().Unix(),
			}
		}
	}

	return nil
}

func (idx index) write(dir string) error {

	// don't write index contents as is - this could lead to data loss
	// imagine having fewer entries in the current index than existing index
	// if we write current index in that case we'll lose all the entries that
	// are present in existing index, but not current
	//
	// to avoid that - read existing index in the same dir first
	// then replace values for each key from the current index
	// finally write that updated index to disk
	//
	// this would avoid data loss, if the index was loaded for fewer values
	// than existing index contains

	exind := make(index)
	if err := exind.read(dir); err != nil {
		return err
	}

	for k, v := range idx {
		exind[k] = v
	}

	indexPath := indexPath(dir)

	indexFile, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	return gob.NewEncoder(indexFile).Encode(exind)
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
