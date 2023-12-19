package kvas

import "github.com/boggydigital/nod"

func ReduxVetter(dir string, assets ...string) (IndexVetter, error) {
	return connectRedux(dir, assets...)
}

func (rdx *redux) VetIndexOnly(fix bool, tpw nod.TotalProgressWriter) ([]string, error) {
	return rdx.kv.VetIndexOnly(fix, tpw)
}

func (rdx *redux) VetIndexMissing(fix bool, tpw nod.TotalProgressWriter) ([]string, error) {
	return rdx.kv.VetIndexMissing(fix, tpw)
}
