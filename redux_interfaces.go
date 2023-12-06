package kvas

type ValueEditor interface {
	AddValues(key string, values ...string) error
	ReplaceValues(key string, values ...string) error
	CutVal(key string, val string) error
}

type BatchValueEditor interface {
	BatchAddValues(keysValues map[string][]string) error
	BatchReplaceValues(keysValues map[string][]string) error
	BatchCutValues(keysValues map[string][]string) error
}

type ValueChecker interface {
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
	ValueChecker
	AllValuesGetter
	FirstValueGetter
	TermsMatcher
}

type ReduxRefresher interface {
	RefreshReduxValues() (ReduxValues, error)
	ReduxModTime() (int64, error)
}

type ReduxValues interface {
	ValueEditor
	BatchValueEditor
	ValueReader
	ReduxRefresher
}
