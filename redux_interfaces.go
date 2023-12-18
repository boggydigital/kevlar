package kvas

import "io"

type ReadableRedux interface {
	MustHave(assets ...string) error
	Keys(asset string) []string
	Has(asset string) bool
	HasKey(asset, key string) bool
	HasVal(asset, key, val string) bool
	GetFirstVal(asset, key string) (string, bool)
	GetAllValues(asset, key string) ([]string, bool)
	ModTime() (int64, error)
	Refresh() error
	MatchAsset(asset string, terms []string, scope []string, anyCase bool, contains bool) []string
	Match(query map[string][]string, anyCase, contains bool) []string
	Sort(ids []string, desc bool, sortBy ...string) ([]string, error)
	Export(w io.Writer, ids ...string) error
}

type WriteableRedux interface {
	ReadableRedux
	AddValues(asset, key string, values ...string) error
	BatchAddValues(asset string, keyValues map[string][]string) error
	ReplaceValues(asset, key string, values ...string) error
	BatchReplaceValues(asset string, keyValues map[string][]string) error
	CutValues(asset, key string, values ...string) error
	BatchCutValues(asset string, keyValues map[string][]string) error
	BatchCutKeys(asset string, keys []string) error
}

type ToolableRedux interface {
}
