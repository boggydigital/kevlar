package kvas

import "io"

type KeyValuesRefresher interface {
	IndexCurrentModTime() (int64, error)
	CurrentModTime(key string) (int64, error)
	IndexRefresh() error
}

type KeyValues interface {
	KeyValuesRefresher

	Has(key string) bool
	Get(key string) (io.ReadCloser, error)
	GetFromStorage(key string) (io.ReadCloser, error)
	Set(key string, data io.Reader) error
	Cut(key string) (bool, error)

	Keys() []string
	CreatedAfter(timestamp int64) []string
	ModifiedAfter(timestamp int64, strictlyModified bool) []string
	IsModifiedAfter(key string, timestamp int64) bool
}
