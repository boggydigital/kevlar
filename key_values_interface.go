package kevlar

import (
	"io"
	"time"
)

type KeyValues interface {
	Keys() ([]string, error)
	Has(key string) (bool, error)

	Get(key string) (io.ReadCloser, error)
	Set(key string, data io.Reader) error
	Cut(key string) (bool, error)

	IsCurrent() (bool, time.Time)
	CreatedAfter(ts int64) ([]string, error)
	UpdatedAfter(ts int64) ([]string, error)
	CreatedOrUpdatedAfter(ts int64) ([]string, error)
	IsUpdatedAfter(key string, ts int64) (bool, error)

	ModTime(key string) (time.Time, error)
}
