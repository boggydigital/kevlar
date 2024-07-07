package kevlar

import "time"

func (rdx *redux) ModTime() (time.Time, error) {
	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	mt := time.Unix(0, 0)
	for asset := range rdx.akv {
		if amt, err := rdx.kv.ModTime(asset); err != nil {
			return time.Unix(0, 0), err
		} else {
			if amt.After(mt) {
				mt = amt
			}
		}
	}
	return mt, nil
}

func (rdx *redux) refresh() (*redux, error) {

	for asset := range rdx.akv {
		if ok, _ := rdx.kv.IsCurrent(); ok {
			continue
		}

		ckv, err := loadAsset(rdx.kv, asset)
		if err != nil {
			return nil, err
		}
		rdx.akv[asset] = ckv
	}

	return rdx, nil
}

func (rdx *redux) RefreshReader() (ReadableRedux, error) {
	return rdx.refresh()
}
