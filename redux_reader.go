package kvas

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func ReduxReader(dir string, assets ...string) (ReadableRedux, error) {
	return connectRedux(dir, assets...)
}

func (rdx *Redux) MustHave(assets ...string) error {
	for _, asset := range assets {
		if !rdx.HasAsset(asset) {
			return UnknownReduxAsset(asset)
		}
	}
	return nil
}

func (rdx *Redux) Keys(asset string) []string {
	return maps.Keys(rdx.assetKeyValues[asset])
}

func (rdx *Redux) HasAsset(asset string) bool {
	_, ok := rdx.assetKeyValues[asset]
	return ok
}

func (rdx *Redux) HasKey(asset, key string) bool {
	if akr, ok := rdx.assetKeyValues[asset]; ok {
		_, ok = akr[key]
		return ok
	} else {
		return false
	}
}

func (rdx *Redux) HasValue(asset, key, val string) bool {
	if akr, ok := rdx.assetKeyValues[asset]; ok {
		if kr, ok := akr[key]; ok {
			return slices.Contains(kr, val)
		} else {
			return false
		}
	} else {
		return false
	}
}

func (rdx *Redux) GetAllValues(asset, key string) ([]string, bool) {
	if !rdx.HasAsset(asset) {
		return nil, false
	}
	if rdx.assetKeyValues[asset] == nil {
		return nil, false
	}

	val, ok := rdx.assetKeyValues[asset][key]
	return val, ok
}

func (rdx *Redux) GetFirstVal(asset, key string) (string, bool) {
	if values, ok := rdx.GetAllValues(asset, key); ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}
