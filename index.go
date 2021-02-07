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
	ip := filepath.Join(vs.baseDir, indexFilename+vs.ext)
	if vs.debug {
		log.Println("ValueSet.indexPath:", ip)
	}
	return ip
}

// initIndex initializes index data structure
func (vs *ValueSet) initIndex() {
	if vs.index == nil {
		if vs.debug {
			log.Println("ValueSet.initIndex: initialized nil index")
		}
		vs.index = make(map[string]IndexRecord, 0)
	}
}

// readIndex reads index of a valueSet
func (vs *ValueSet) readIndex() error {
	indexPath := vs.indexPath()

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		if vs.debug {
			log.Println("ValueSet.readIndex: index file doesn't exist at", indexPath)
		}
		return nil
	}

	bytes, err := ioutil.ReadFile(indexPath)
	if err != nil {
		if vs.debug {
			log.Println("ValueSet.readIndex: error reading index file", err.Error())
		}
		return err
	}

	if vs.debug {
		log.Printf("ValueSet.readIndex: read %d bytes", len(bytes))
	}

	return json.Unmarshal(bytes, &vs.index)
}

// writeIndex writes index of a valueSet
func (vs *ValueSet) writeIndex() error {
	bytes, err := json.Marshal(vs.index)
	if err != nil {
		if vs.debug {
			log.Println("ValueSet.writeIndex: error marshalling index", err.Error())
		}
		return err
	}

	if vs.debug {
		log.Printf("ValueSet.writeIndex: marshalled index to %d bytes\n", len(bytes))
	}

	return ioutil.WriteFile(vs.indexPath(), bytes, filePerm)
}

// setIndex updates index by key
func (vs *ValueSet) setIndex(key string, hash string) {
	if _, ok := vs.index[key]; !ok {
		if vs.debug {
			log.Printf("ValueSet.setIndex: no entry for the key '%s'\n", key)
		}
		vs.index[key] = IndexRecord{
			Created: time.Now().Unix(),
		}
	}

	if vs.index[key].Hash != hash {
		if vs.debug {
			log.Printf("ValueSet.setIndex: hash for item with key '%s' has changed\n", key)
		}
		vs.index[key] = IndexRecord{
			Hash:     hash,
			Created:  vs.index[key].Created,
			Modified: time.Now().Unix(),
		}
	} else {
		log.Printf("ValueSet.setIndex: hash for item with key '%s' is the same", key)
	}
}
