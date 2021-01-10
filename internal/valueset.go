package internal

type IndexRecord struct {
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type valueSet struct {
	baseDir string
	index   map[string]IndexRecord
}

// ValueSetClient defines functions of a key value store client
type ValueSetClient interface {
	getExt() string
	// index
	indexPath() string
	initIndex()
	readIndex() error
	writeIndex() error
	removeIndex(string)
	setIndex(string, string)
	// values
	valuePath(string) string
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Remove(string) error
	Contains(string) bool
	// enumerations
	reduce(func(IndexRecord) bool) []string
	All() []string
	CreatedAfter(int64) []string
	ModifiedAfter(int64) []string
}

// getExt returns file extension for the files in a valueSet
func (vs *valueSet) getExt() string {
	return ""
}

// Contains verifies if a value set contains provided key
func (vs *valueSet) Contains(key string) bool {
	_, ok := vs.index[key]
	return ok
}
