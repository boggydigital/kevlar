package kvas

import "golang.org/x/exp/maps"

func (rdx *redux) assetModTime(asset string) (int64, error) {
	return rdx.kv.CurrentModTime(asset)
}

func (rdx *redux) ModTime() (int64, error) {
	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	var mt int64 = 0
	for asset := range rdx.assetKeyValues {
		if amt, err := rdx.assetModTime(asset); err != nil {
			return -1, err
		} else {
			if mt < amt {
				mt = amt
			}
		}
	}
	return mt, nil
}

func (rdx *redux) refresh() (*redux, error) {
	if err := rdx.kv.IndexRefresh(); err != nil {
		return rdx, err
	}

	modTime, err := rdx.ModTime()
	if err != nil {
		return rdx, err
	}

	if rdx.modTime < modTime {
		return connectRedux(rdx.dir, maps.Keys(rdx.assetKeyValues)...)
	}

	return rdx, nil
}

func (rdx *redux) RefreshReader() (ReadableRedux, error) {
	return rdx.refresh()
}
