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
	Cut(key string) error

	Since(ts int64, mts ...MutationType) iter.Seq2[string, MutationType]
	//CreatedAfter(ts int64) iter.Seq[string]
	//UpdatedAfter(ts int64) iter.Seq[string]
	//CreatedOrUpdatedAfter(ts int64) iter.Seq[string]

	LogModTime(key string) int64
	FileModTime(key string) (int64, error)
}
