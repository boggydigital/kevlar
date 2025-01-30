package kevlar

import (
	"io"
	"iter"
)

type KeyValues interface {
	Len() int
	Keys() iter.Seq[string]
	Has(key string) bool

	Get(key string) (io.ReadCloser, error)
	Set(key string, data io.Reader) error
	Cut(key string) (bool, error)

	CreatedAfter(ts int64) iter.Seq[string]
	UpdatedAfter(ts int64) iter.Seq[string]
	CreatedOrUpdatedAfter(ts int64) iter.Seq[string]
	IsUpdatedAfter(key string, ts int64) bool

	ModTime() int64
	ValueModTime(key string) int64
}
