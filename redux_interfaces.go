package kvas

type ValueAdder interface {
	AddVal(string, string) error
}

type ValuesReplacer interface {
	ReplaceValues(string, ...string) error
}

type BatchValuesReplacer interface {
	BatchReplaceValues(map[string][]string) error
}

type ValueCutter interface {
	CutVal(string, string) error
}

type ValueEditor interface {
	ValueAdder
	ValuesReplacer
	BatchValuesReplacer
	ValueCutter
}

type ValuePresenceChecker interface {
	HasVal(string, string) bool
}

type AllValuesGetter interface {
	GetAllValues(string) ([]string, bool)
}

type FirstValueGetter interface {
	GetFirstVal(string) (string, bool)
}

type TermsMatcher interface {
	Match([]string, map[string]bool, bool, bool) map[string]bool
}

type ValueReader interface {
	KeysEnumerator
	PresenceChecker
	ValuePresenceChecker
	AllValuesGetter
	FirstValueGetter
	TermsMatcher
}

type ValueRefresher interface {
	RefreshReduxValues() (ReduxValues, error)
}

type ReduxValues interface {
	ValueEditor
	ValueReader
	ValueRefresher
}
