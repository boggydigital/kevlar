package kvas

import "io"

type Container interface {
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

type KeyValuesClient interface {
	Container
	Getter
	Setter
	Cutter
}

type KeyValuesFilters interface {
	KeyValuesClient
	CreatedAfter(int64) []string
	ModifiedAfter(int64, bool) []string
	IsModifiedAfter(string, int64) bool
}
