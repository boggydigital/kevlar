package kevlar

import (
	"encoding/gob"
	"errors"
	"io"
	"sync"
)

func ErrUnknownAsset(asset string) error {
	return errors.New("unknown redux asset " + asset)
}

type redux struct {
	dir string
	kv  KeyValues
	akv map[string]map[string][]string
	lmt map[string]int64
	mtx *sync.Mutex
}

func newRedux(dir string, assets ...string) (*redux, error) {
	kv, err := NewKeyValues(dir, GobExt)
	if err != nil {
		return nil, err
	}

	assetKeyValues := make(map[string]map[string][]string)
	amts := make(map[string]int64)
	for _, asset := range assets {
		if assetKeyValues[asset], err = loadAsset(kv, asset); err != nil {
			return nil, err
		}
		amts[asset], err = kv.ModTime(asset)
		if err != nil {
			return nil, err
		}
	}

	return &redux{
		kv:  kv,
		dir: dir,
		akv: assetKeyValues,
		lmt: amts,
		mtx: new(sync.Mutex),
	}, nil
}

func loadAsset(kv KeyValues, asset string) (map[string][]string, error) {

	ok, err := kv.Has(asset)
	if err != nil {
		return nil, err
	}
	if !ok {
		return make(map[string][]string), nil
	}

	arc, err := kv.Get(asset)
	if err != nil {
		return nil, err
	}
	if arc != nil {
		defer arc.Close()
	}

	var keyValues map[string][]string
	if arc != nil {
		if err := gob.NewDecoder(arc).Decode(&keyValues); err == io.EOF {
			// empty reduction - do nothing, it'll be initialized below
		} else if err != nil {
			return nil, err
		}
	}

	if keyValues == nil {
		keyValues = make(map[string][]string)
	}

	return keyValues, nil
}
