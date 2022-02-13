package kvas

//ReduxAtomics enumerates assets that shouldn't be matched by containment, only as a whole word.
type ReduxAtomics map[string]bool

func (ra ReduxAtomics) IsAtomic(asset string) bool {
	_, ok := ra[asset]
	return ok
}
