package kevlar

// assetModTimes returns ModTime for each asset. It doesn't Update it
// because that time should be updated only when asset is loaded
func (rdx *redux) assetsModTimes() map[string]int64 {
	amts := make(map[string]int64)
	for asset := range rdx.akv {
		amts[asset] = rdx.kv.ValueModTime(asset)
	}
	return amts
}

func (rdx *redux) ModTime() int64 {
	var mt int64 = UnknownModTime
	amts := rdx.assetsModTimes()

	for asset := range rdx.akv {
		if amt, ok := amts[asset]; ok && amt > mt {
			mt = amt
		}
	}

	return mt
}

func (rdx *redux) refresh() (*redux, error) {

	amts := rdx.assetsModTimes()
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
