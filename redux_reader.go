package kvas

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func NewReduxReader(dir string, assets ...string) (ReadableRedux, error) {
	return connectRedux(dir, assets...)
}

func (rdx *redux) MustHave(assets ...string) error {
	for _, asset := range assets {
		if !rdx.HasAsset(asset) {
			return UnknownReduxAsset(asset)
		}
	}
	return nil
}

func (rdx *redux) Keys(asset string) []string {
	return maps.Keys(rdx.assetKeyValues[asset])
}

func (rdx *redux) HasAsset(asset string) bool {
	_, ok := rdx.assetKeyValues[asset]
	return ok
}

func (rdx *redux) HasKey(asset, key string) bool {
	if akr, ok := rdx.assetKeyValues[asset]; ok {
		_, ok = akr[key]
		return ok
	} else {
		return false
	}
}

func (rdx *redux) HasValue(asset, key, val string) bool {
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

func (rdx *redux) GetAllValues(asset, key string) ([]string, bool) {
	if !rdx.HasAsset(asset) {
		return nil, false
	}
	if rdx.assetKeyValues[asset] == nil {
		return nil, false
	}

	val, ok := rdx.assetKeyValues[asset][key]
	return val, ok
}

func (rdx *redux) GetFirstVal(asset, key string) (string, bool) {
	if values, ok := rdx.GetAllValues(asset, key); ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (rdx *redux) GetLastVal(asset, key string) (string, bool) {
	if values, ok := rdx.GetAllValues(asset, key); ok && len(values) > 0 {
		return values[len(values)-1], true
	}
	return "", false
}
