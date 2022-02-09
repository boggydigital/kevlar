package kvas

type AssetKeysEnumerator interface {
	Keys(string) []string
}

type AssetPresenceChecker interface {
	Has(string) bool
}

type AssetKeyPresenceChecker interface {
	HasKey(string, string) bool
}

type AssetValuePresenceChecker interface {
	HasVal(string, string, string) bool
}

type AssetValueAdder interface {
	AddVal(string, string, string) error
}

type AssetValuesReplacer interface {
	ReplaceValues(string, string, ...string) error
}

type AssetBatchValuesReplacer interface {
	BatchReplaceValues(string, map[string][]string) error
}

type AssetValueCutter interface {
	CutVal(string, string, string) error
}

type AssetAllValuesGetter interface {
	GetAllValues(string, string) ([]string, bool)
}

type AssetAllUnchangedValuesGetter interface {
	GetAllUnchangedValues(string, string) ([]string, bool)
}

type AssetFirstValueGetter interface {
	GetFirstVal(string, string) (string, bool)
}

type AssetsSupportChecker interface {
	IsSupported(...string) error
}

type QueryMatcher interface {
	Match(map[string][]string, bool) map[string]bool
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
}
