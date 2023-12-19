package kvas

import (
	"github.com/boggydigital/nod"
	"io"
)

type IndexVetter interface {
	VetIndexOnly(fix bool, tpw nod.TotalProgressWriter) ([]string, error)
	VetIndexMissing(fix bool, tpw nod.TotalProgressWriter) ([]string, error)
}

type KeyValues interface {
	Has(key string) bool
	Get(key string) (io.ReadCloser, error)
	GetFromStorage(key string) (io.ReadCloser, error)
	Set(key string, data io.Reader) error
	Cut(key string) (bool, error)

	Keys() []string
	CreatedAfter(timestamp int64) []string
	ModifiedAfter(timestamp int64, strictlyModified bool) []string
	IsModifiedAfter(key string, timestamp int64) bool

	IndexCurrentModTime() (int64, error)
	CurrentModTime(key string) (int64, error)
	IndexRefresh() error

	IndexVetter
}
