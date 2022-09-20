package kvas

type AssetKeysEnumerator interface {
	Keys(asset string) []string
}

type AssetPresenceChecker interface {
	Has(asset string) bool
}

type AssetKeyPresenceChecker interface {
	HasKey(asset, key string) bool
}

type AssetValuePresenceChecker interface {
	HasVal(asset, key, val string) bool
}

type AssetValueAdder interface {
	AddVal(asset, key, val string) error
}

type AssetValuesReplacer interface {
	ReplaceValues(asset, key string, values ...string) error
}

type AssetBatchValuesReplacer interface {
	BatchReplaceValues(asset string, keyValues map[string][]string) error
}

type AssetValueCutter interface {
	CutVal(asset, key, val string) error
}

type AssetAllValuesGetter interface {
	GetAllValues(asset, key string) ([]string, bool)
}

type AssetAllUnchangedValuesGetter interface {
	GetAllUnchangedValues(asset, key string) ([]string, bool)
}

type AssetFirstValueGetter interface {
	GetFirstVal(asset, key string) (string, bool)
}

type AssetsSupportChecker interface {
	IsSupported(assets ...string) error
}

type QueryMatcher interface {
	Match(query map[string][]string, anyCase bool) map[string]bool
}

type AssetsRefresher interface {
	RefreshReduxAssets() (ReduxAssets, error)
}

type AssetsModTimeGetter interface {
	ReduxAssetsModTime() (int64, error)
}

type AssetsSorter interface {
	Sort(ids []string, desc bool, sortBy ...string) ([]string, error)
}

type ReduxAssets interface {
	AssetKeysEnumerator
	AssetPresenceChecker
	AssetKeyPresenceChecker
	AssetValuePresenceChecker
	AssetValueAdder
	AssetValuesReplacer
	AssetBatchValuesReplacer
	AssetValueCutter
	AssetAllValuesGetter
	AssetAllUnchangedValuesGetter
	AssetFirstValueGetter
	AssetsSupportChecker
	QueryMatcher
	AssetsRefresher
	AssetsModTimeGetter
	AssetsSorter
}
