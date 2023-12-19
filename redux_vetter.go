package kvas

func ReduxVetter(dir string, assets ...string) (IndexVetter, error) {
	return connectRedux(dir, assets...)
}

func (rdx *redux) VetIndexOnly(fix bool) ([]string, error) {
	return rdx.kv.VetIndexOnly(fix)
}

func (rdx *redux) VetIndexMissing(fix bool) ([]string, error) {
	return rdx.kv.VetIndexMissing(fix)
}
