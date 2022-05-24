package kvas

type ValueAdder interface {
	AddVal(key string, val string) error
}

type ValuesReplacer interface {
	ReplaceValues(key string, values ...string) error
}

type BatchValuesReplacer interface {
	BatchReplaceValues(keysValues map[string][]string) error
}

type ValueCutter interface {
	CutVal(key string, val string) error
}

type ValueEditor interface {
	ValueAdder
	ValuesReplacer
	BatchValuesReplacer
	ValueCutter
}

type ValuePresenceChecker interface {
	HasVal(key string, val string) bool
}

type AllValuesGetter interface {
	GetAllValues(key string) ([]string, bool)
}

type FirstValueGetter interface {
	GetFirstVal(key string) (string, bool)
}

type TermsMatcher interface {
	Match(terms []string, scope map[string]bool, anyCase bool, contains bool) map[string]bool
}

type ValueReader interface {
	KeysEnumerator
	PresenceChecker
	ValuePresenceChecker
	AllValuesGetter
	FirstValueGetter
	TermsMatcher
}

type ReduxRefresher interface {
	RefreshReduxValues() (ReduxValues, error)
}

type ReduxModTimeGetter interface {
	ReduxModTime() (int64, error)
}

type ReduxValues interface {
	ValueEditor
	ValueReader
	ReduxRefresher
	ReduxModTimeGetter
}
