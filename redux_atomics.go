package kvas

//reduxAtomics enumerates assets that shouldn't be matched by containment, only as a whole word.
type reduxAtomics map[string]bool

func (ra reduxAtomics) IsAtomic(asset string) bool {
	_, ok := ra[asset]
	return ok
}
