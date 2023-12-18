package kvas

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/exp/slices"
)

func (rdx *Redux) AddValues(asset, key string, values ...string) error {
	if !rdx.Has(asset) {
		return UnknownReduxAsset(asset)
	}
	newValues := make([]string, 0, len(values))
	for _, v := range values {
		if !rdx.HasValue(asset, key, v) {
			newValues = append(newValues, v)
		}
	}
	rdx.assetKeyValues[asset][key] = append(rdx.assetKeyValues[asset][key], newValues...)
	return rdx.write(asset)
}

func (rdx *Redux) BatchAddValues(asset string, keyValues map[string][]string) error {
	if !rdx.Has(asset) {
		return UnknownReduxAsset(asset)
	}
	if len(keyValues) == 0 {
		return nil
	}
	for key, values := range keyValues {
		for _, v := range values {
			if !rdx.HasValue(asset, key, v) {
				rdx.assetKeyValues[asset][key] = append(rdx.assetKeyValues[asset][key], v)
			}
		}
	}
	return rdx.write(asset)
}

func (rdx *Redux) replaceValues(asset, key string, values ...string) error {
	if !rdx.Has(asset) {
		return UnknownReduxAsset(asset)
	}
	rdx.assetKeyValues[asset][key] = values
	return nil
}

func (rdx *Redux) ReplaceValues(asset, key string, values ...string) error {
	if err := rdx.replaceValues(asset, key, values...); err != nil {
		return err
	}
	return rdx.write(asset)
}

func (rdx *Redux) BatchReplaceValues(asset string, keyValues map[string][]string) error {
	if len(keyValues) == 0 {
		return nil
	}
	for key, values := range keyValues {
		if err := rdx.replaceValues(asset, key, values...); err != nil {
			return err
		}
	}
	return rdx.write(asset)
}

func (rdx *Redux) cutValues(asset, key string, values ...string) error {
	if !rdx.Has(asset) {
		return UnknownReduxAsset(asset)
	}
	if !rdx.HasKey(asset, key) {
		return nil
	}

	newValues := make([]string, 0, len(rdx.assetKeyValues[asset][key]))

	for _, v := range rdx.assetKeyValues[asset][key] {
		if slices.Contains(values, v) {
			continue
		}
		newValues = append(newValues, v)
	}

	rdx.assetKeyValues[asset][key] = newValues

	// remove keys if there are no values left
	if len(rdx.assetKeyValues[asset][key]) == 0 {
		delete(rdx.assetKeyValues[asset], key)
	}
	return nil
}

func (rdx *Redux) CutValues(asset, key string, values ...string) error {
	if err := rdx.cutValues(asset, key, values...); err != nil {
		return err
	}
	return rdx.write(asset)
}

func (rdx *Redux) BatchCutValues(asset string, keyValues map[string][]string) error {
	if len(keyValues) == 0 {
		return nil
	}
	for key, values := range keyValues {
		if err := rdx.cutValues(asset, key, values...); err != nil {
			return err
		}
	}
	return rdx.write(asset)
}

func (rdx *Redux) BatchCutKeys(asset string, keys []string) error {
	if !rdx.Has(asset) {
		return UnknownReduxAsset(asset)
	}
	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		delete(rdx.assetKeyValues[asset], key)
	}
	return rdx.write(asset)
}

func (rdx *Redux) write(asset string) error {
	if !rdx.Has(asset) {
		return UnknownReduxAsset(asset)
	}

	kv, err := ConnectLocal(rdx.dir, GobExt)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(rdx.assetKeyValues[asset]); err != nil {
		return err
	}

	return kv.Set(asset, buf)
}
