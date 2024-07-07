package kevlar

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/exp/slices"
)

func NewReduxWriter(dir string, assets ...string) (WriteableRedux, error) {
	return newRedux(dir, assets...)
}

func (rdx *redux) addValues(asset, key string, values ...string) error {
	if !rdx.HasAsset(asset) {
		return ErrUnknownAsset(asset)
	}
	newValues := make([]string, 0, len(values))
	for _, v := range values {
		if !rdx.HasValue(asset, key, v) {
			newValues = append(newValues, v)
		}
	}
	rdx.akv[asset][key] = append(rdx.akv[asset][key], newValues...)
	return rdx.write(asset)
}

func (rdx *redux) AddValues(asset, key string, values ...string) error {
	return rdx.addValues(asset, key, values...)
}

func (rdx *redux) BatchAddValues(asset string, keyValues map[string][]string) error {
	for key, values := range keyValues {
		if err := rdx.addValues(asset, key, values...); err != nil {
			return err
		}
	}
	return nil
}

func (rdx *redux) replaceValues(asset, key string, values ...string) error {
	if !rdx.HasAsset(asset) {
		return ErrUnknownAsset(asset)
	}
	rdx.akv[asset][key] = values
	return nil
}

func (rdx *redux) ReplaceValues(asset, key string, values ...string) error {
	if err := rdx.replaceValues(asset, key, values...); err != nil {
		return err
	}
	return rdx.write(asset)
}

func (rdx *redux) BatchReplaceValues(asset string, keyValues map[string][]string) error {
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

func (rdx *redux) cutValues(asset, key string, values ...string) error {
	if !rdx.HasAsset(asset) {
		return ErrUnknownAsset(asset)
	}
	if !rdx.HasKey(asset, key) {
		return nil
	}

	newValues := make([]string, 0, len(rdx.akv[asset][key]))

	for _, v := range rdx.akv[asset][key] {
		if slices.Contains(values, v) {
			continue
		}
		newValues = append(newValues, v)
	}

	rdx.akv[asset][key] = newValues

	// remove keys if there are no values left
	if len(rdx.akv[asset][key]) == 0 {
		delete(rdx.akv[asset], key)
	}
	return nil
}

func (rdx *redux) CutValues(asset, key string, values ...string) error {
	if err := rdx.cutValues(asset, key, values...); err != nil {
		return err
	}
	return rdx.write(asset)
}

func (rdx *redux) CutKeys(asset string, keys ...string) error {
	if !rdx.HasAsset(asset) {
		return ErrUnknownAsset(asset)
	}
	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		delete(rdx.akv[asset], key)
	}
	return rdx.write(asset)
}

func (rdx *redux) BatchCutValues(asset string, keyValues map[string][]string) error {
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

func (rdx *redux) write(asset string) error {
	if !rdx.HasAsset(asset) {
		return ErrUnknownAsset(asset)
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(rdx.akv[asset]); err != nil {
		return err
	}

	return rdx.kv.Set(asset, buf)
}

func (rdx *redux) RefreshWriter() (WriteableRedux, error) {
	return rdx.refresh()
}
