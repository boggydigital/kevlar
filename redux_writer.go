package kvas

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/exp/slices"
)

func ReduxWriter(dir string, assets ...string) (WriteableRedux, error) {
	return connectRedux(dir, assets...)
}

func (rdx *Redux) addValues(asset, key string, values ...string) error {
	if !rdx.HasAsset(asset) {
		return UnknownReduxAsset(asset)
	}
	newValues := make([]string, 0, len(values))
	for _, v := range values {
		if !rdx.HasValue(asset, key, v) {
			newValues = append(newValues, v)
		}
	}
	rdx.assetKeyValues[asset][key] = append(rdx.assetKeyValues[asset][key], newValues...)
	return nil
}

func (rdx *Redux) AddValues(asset, key string, values ...string) error {
	if err := rdx.addValues(asset, key, values...); err != nil {
		return err
	}
	return rdx.write(asset)
}

func (rdx *Redux) BatchAddValues(asset string, keyValues map[string][]string) error {
	for key, values := range keyValues {
		if err := rdx.addValues(asset, key, values...); err != nil {
			return err
		}
	}
	return rdx.write(asset)
}

func (rdx *Redux) replaceValues(asset, key string, values ...string) error {
	if !rdx.HasAsset(asset) {
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
	if !rdx.HasAsset(asset) {
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
	if !rdx.HasAsset(asset) {
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
	if !rdx.HasAsset(asset) {
		return UnknownReduxAsset(asset)
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(rdx.assetKeyValues[asset]); err != nil {
		return err
	}

	return rdx.kv.Set(asset, buf)
}

func (rdx *Redux) RefreshWriter() (WriteableRedux, error) {
	return rdx.refresh()
}
