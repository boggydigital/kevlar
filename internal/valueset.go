package internal

import "time"

const (
	indexFilename = "_index"
)

type IndexRecord struct {
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type valueSet struct {
	baseDir string
	index   map[string]IndexRecord
}

type ValueSetClient interface {
	indexPath() string
	valuePath(string) string
	initIndex()
	readIndex() error
	writeIndex() error
	removeIndex(string)
	setIndex(string, string)
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Remove(string) error
	Contains(string) bool
	All() []string
	CreatedAfter(int64) []string
	ModifiedAfter(int64) []string
}

func (vs *valueSet) initIndex() {
	if vs.index == nil {
		vs.index = make(map[string]IndexRecord, 0)
	}
}

func (vs *valueSet) setIndex(key string, hash string) {
	if _, ok := vs.index[key]; !ok {
		vs.index[key] = IndexRecord{
			Created: time.Now().Unix(),
		}
	}

	vs.index[key] = IndexRecord{
		Hash:     hash,
		Created:  vs.index[key].Created,
		Modified: time.Now().Unix(),
	}
}

func (vs *valueSet) Contains(key string) bool {
	if _, ok := vs.index[key]; ok {
		return true
	}
	return false
}

func (vs *valueSet) All() []string {
	if vs == nil {
		return nil
	}
	keys := make([]string, 0, len(vs.index))
	for k := range vs.index {
		keys = append(keys, k)
	}
	return keys
}

func (vs *valueSet) CreatedAfter(timestamp int64) []string {
	if vs == nil {
		return nil
	}
	keys := make([]string, 0, len(vs.index))
	for k, ir := range vs.index {
		if ir.Created > timestamp {
			keys = append(keys, k)
		}
	}
	return keys
}

func (vs *valueSet) ModifiedAfter(timestamp int64) []string {
	if vs == nil {
		return nil
	}
	keys := make([]string, 0, len(vs.index))
	for k, ir := range vs.index {
		if ir.Modified > timestamp {
			keys = append(keys, k)
		}
	}
	return keys
}
