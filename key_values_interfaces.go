package kvas

import "io"

type PresenceChecker interface {
	Has(string) bool
}

type Getter interface {
	Get(string) (io.ReadCloser, error)
}

type Setter interface {
	Set(string, io.Reader) error
}

type Cutter interface {
	Cut(string) (bool, error)
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
	CreatedAfter(int64) []string
}

type ModifiedAfterFilter interface {
	ModifiedAfter(int64, bool) []string
}

type ModifiedAfterChecker interface {
	IsModifiedAfter(string, int64) bool
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
