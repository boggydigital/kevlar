package internal

import (
	"time"
)

// initIndex initializes index data structure
func (vs *valueSet) initIndex() {
	if vs.index == nil {
		vs.index = make(map[string]IndexRecord, 0)
	}
}

// setIndex updates index by key
func (vs *valueSet) setIndex(key string, hash string) {
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
	}
}
