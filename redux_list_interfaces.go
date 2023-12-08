package kvas

import "io"

type AssetReader interface {
	Keys(asset string) []string
	GetFirstVal(asset, key string) (string, bool)
	IsSupported(assets ...string) error
	GetAllValues(asset, key string) ([]string, bool)
}

type AssetChecker interface {
	Has(asset string) bool
	HasKey(asset, key string) bool
	HasVal(asset, key, val string) bool
	MustHave(asset string) error
}

type AssetEditor interface {
	AddValues(asset, key string, values ...string) error
	ReplaceValues(asset, key string, values ...string) error
	CutVal(asset, key, val string) error
}

type AssetBatchEditor interface {
	BatchAddValues(asset string, keyValues map[string][]string) error
	BatchReplaceValues(asset string, keyValues map[string][]string) error
	BatchCutKeys(asset string, keys []string) error
	BatchCutValues(asset string, keyValues map[string][]string) error
}

type QueryMatcher interface {
	Match(query map[string][]string, anyCase, contains bool) map[string]bool
}

type AssetsRefresher interface {
	RefreshReduxAssets() (ReduxAssets, error)
	ReduxAssetsModTime() (int64, error)
}

type AssetsSorter interface {
	Sort(ids []string, desc bool, sortBy ...string) ([]string, error)
}

type AssetsExporter interface {
	Export(w io.Writer, ids ...string) error
}

type ReduxAssets interface {
	AssetChecker
	AssetReader
	AssetEditor
	AssetBatchEditor
	QueryMatcher
	AssetsRefresher
	AssetsSorter
	AssetsExporter
}
