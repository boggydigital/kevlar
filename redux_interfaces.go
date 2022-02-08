package kvas

type ValuePresenceChecker interface {
	HasVal(string, string) bool
}

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

type AllValuesGetter interface {
	GetAllValues(string) ([]string, bool)
}

type FirstValueGetter interface {
	GetFirstVal(string) (string, bool)
}

type TermsMatcher interface {
	Match([]string, []string, bool, bool) map[string]bool
}

type ReduxValues interface {
	ValuePresenceChecker
	ValueAdder
	ValuesReplacer
	BatchValuesReplacer
	ValueCutter
	AllValuesGetter
	FirstValueGetter
	TermsMatcher
}
