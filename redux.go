package kvas

import (
	"encoding/gob"
	"errors"
	"io"
	"sync"
	"time"
)

func UnknownReduxAsset(asset string) error {
	return errors.New("unknown redux asset " + asset)
}

type redux struct {
	dir            string
	kv             KeyValues
	assetKeyValues map[string]map[string][]string
	modTime        int64
	mtx            *sync.Mutex
}

func connectRedux(dir string, assets ...string) (*redux, error) {
	kv, err := NewKeyValues(dir, GobExt)
	if err != nil {
		return nil, err
	}

	assetKeyValues := make(map[string]map[string][]string)
	for _, asset := range assets {
		if assetKeyValues[asset], err = loadAsset(kv, asset); err != nil {
			return nil, err
		}
	}

	return &redux{
		kv:             kv,
		dir:            dir,
		assetKeyValues: assetKeyValues,
		modTime:        time.Now().Unix(),
		mtx:            &sync.Mutex{},
	}, nil
}

func loadAsset(kvr KeyValues, asset string) (map[string][]string, error) {
	arc, err := kvr.Get(asset)
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
