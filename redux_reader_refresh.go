package kevlar

// assetsModTimes returns the later of FileModTime, LogModTime for each asset. It doesn't Update it
// because that time should be updated only when asset is loaded
func (rdx *redux) assetsModTimes() (map[string]int64, error) {
	amt := make(map[string]int64)
	var err error
	for asset := range rdx.akv {
		if amt[asset], err = rdx.kv.FileModTime(asset); err != nil {
			return nil, err
		}
		if almt := rdx.kv.LogModTime(asset); almt > amt[asset] {
			amt[asset] = almt
		}
	}
	return amt, nil
}

func (rdx *redux) FileModTime() (int64, error) {
	almt, err := rdx.assetsModTimes()
	if err != nil {
		return UnknownModTime, err
	}

	var mt int64 = UnknownModTime
	for asset := range rdx.akv {
		if amt, ok := almt[asset]; ok && amt > mt {
			mt = amt
		}
	}

	return mt, nil
}

func (rdx *redux) refresh() (*redux, error) {

	amts, err := rdx.assetsModTimes()
	if err != nil {
		return nil, err
	}
	for asset := range rdx.akv {
		// asset was updated externally or not loaded yet
		if rdx.lmt[asset] < amts[asset] {
			ckv, err := loadAsset(rdx.kv, asset)
			if err != nil {
				return nil, err
			}
			rdx.akv[asset] = ckv
			rdx.lmt[asset] = amts[asset]
		}
	}

	return rdx, nil
}

func (rdx *redux) RefreshReader() (ReadableRedux, error) {
	return rdx.refresh()
}
