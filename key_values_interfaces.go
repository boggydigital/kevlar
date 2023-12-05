package kvas

import "io"

type PresenceChecker interface {
	Has(key string) bool
}

type Getter interface {
	Get(key string) (io.ReadCloser, error)
	GetFromStorage(key string) (io.ReadCloser, error)
}

type Setter interface {
	Set(key string, data io.Reader) error
}

type Cutter interface {
	Cut(key string) (bool, error)
}

type KeyValuesEditor interface {
	PresenceChecker
	Getter
	Setter
	Cutter
}

type KeysEnumerator interface {
	Keys() []string
}

type CreatedAfterFilter interface {
	CreatedAfter(timestamp int64) []string
}

type ModifiedAfterFilter interface {
	ModifiedAfter(timestamp int64, strictlyModified bool) []string
}

type ModifiedAfterChecker interface {
	IsModifiedAfter(key string, timestamp int64) bool
}

type KeyValuesFilter interface {
	KeysEnumerator
	CreatedAfterFilter
	ModifiedAfterFilter
	ModifiedAfterChecker
}

type IndexModTimeGetter interface {
	IndexCurrentModTime() (int64, error)
}

type ModTimeGetter interface {
	CurrentModTime(key string) (int64, error)
}

type IndexRefresher interface {
	IndexRefresh() error
}

type KeyValuesRefresher interface {
	IndexModTimeGetter
	ModTimeGetter
	IndexRefresher
}

type KeyValues interface {
	KeyValuesEditor
	KeyValuesFilter
	KeyValuesRefresher
}
