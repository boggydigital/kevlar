package kvas

import (
	"bytes"
	"encoding/gob"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"io"
	"strings"
	"time"
)

type redux struct {
	kvr           KeyValuesRefresher
	dir           string
	asset         string
	keyReductions map[string][]string
	connTime      int64
}

func ConnectRedux(dir, asset string) (ReduxValues, error) {
	rdx, err := ConnectLocal(dir, GobExt)
	if err != nil {
		return nil, err
	}

	arc, err := rdx.Get(asset)
	if arc != nil {
		defer arc.Close()
	}
	if err != nil {
		return nil, err
	}

	ct := time.Now().Unix()

	var keyReductions map[string][]string
	if arc != nil {
		if err := gob.NewDecoder(arc).Decode(&keyReductions); err == io.EOF {
			// empty reduction - do nothing, it'll be initialized below
		} else if err != nil {
			return nil, err
		}
	}

	if keyReductions == nil {
		keyReductions = make(map[string][]string)
	}

	return &redux{
		kvr:           rdx,
		dir:           dir,
		asset:         asset,
		keyReductions: keyReductions,
		connTime:      ct,
	}, nil
}

func (rdx *redux) Keys() []string {
	return maps.Keys(rdx.keyReductions)
}

func (rdx *redux) Has(key string) bool {
	_, ok := rdx.keyReductions[key]
	return ok
}

func (rdx *redux) HasVal(key string, val string) bool {
	for _, v := range rdx.keyReductions[key] {
		if v == val {
			return true
		}
	}
	return false
}

func (rdx *redux) AddValues(key string, values ...string) error {
	newValues := make([]string, 0, len(values))
	for _, val := range values {
		if !rdx.HasVal(key, val) {
			newValues = append(newValues, val)
		}
	}
	rdx.keyReductions[key] = append(rdx.keyReductions[key], newValues...)
	return rdx.write()
}

func (rdx *redux) BatchAddValues(keyValues map[string][]string) error {
	if len(keyValues) == 0 {
		return nil
	}
	for key, values := range keyValues {
		for _, v := range values {
			if slices.Contains(rdx.keyReductions[key], v) {
				continue
			}
			rdx.keyReductions[key] = append(rdx.keyReductions[key], v)
		}
	}
	return rdx.write()
}

func (rdx *redux) ReplaceValues(key string, values ...string) error {
	rdx.keyReductions[key] = values
	return rdx.write()
}

func (rdx *redux) BatchReplaceValues(keysValues map[string][]string) error {
	if len(keysValues) == 0 {
		return nil
	}
	for key, values := range keysValues {
		rdx.keyReductions[key] = values
	}
	return rdx.write()
}

func (rdx *redux) CutVal(key string, val string) error {
	if !rdx.HasVal(key, val) {
		return nil
	}
	values := make([]string, 0, len(rdx.keyReductions[key]))
	for _, v := range rdx.keyReductions[key] {
		if v == val {
			continue
		}
		values = append(values, v)
	}

	rdx.keyReductions[key] = values

	if len(rdx.keyReductions[key]) == 0 {
		delete(rdx.keyReductions, key)
	}
	return rdx.write()
}

func (rdx *redux) BatchCutValues(keyValues map[string][]string) error {
	if len(keyValues) == 0 {
		return nil
	}
	filteredValues := make(map[string][]string)
	for key, values := range rdx.keyReductions {
		filteredValues[key] = make([]string, 0, len(values))
		for _, v := range values {
			if slices.Contains(keyValues[key], v) {
				continue
			}
			filteredValues[key] = append(filteredValues[key], v)
		}
	}
	for key, values := range filteredValues {
		rdx.keyReductions[key] = values
		if len(rdx.keyReductions[key]) == 0 {
			delete(rdx.keyReductions, key)
		}
	}
	return rdx.write()
}

func (rdx *redux) BatchCutKeys(keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		delete(rdx.keyReductions, key)
	}
	return rdx.write()
}

func (rdx *redux) write() error {
	kv, err := ConnectLocal(rdx.dir, GobExt)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(rdx.keyReductions); err != nil {
		return err
	}

	return kv.Set(rdx.asset, buf)
}

func (rdx *redux) GetFirstVal(key string) (string, bool) {
	values, ok := rdx.GetAllValues(key)
	if ok && len(values) > 0 {
		return values[0], ok
	}
	return "", false
}

func (rdx *redux) GetAllValues(key string) ([]string, bool) {
	if rdx == nil || rdx.keyReductions == nil {
		return nil, false
	}
	val, ok := rdx.keyReductions[key]
	return val, ok
}

func (rdx *redux) ReduxModTime() (int64, error) {
	return rdx.kvr.CurrentModTime(rdx.asset)
}

func (rdx *redux) RefreshReduxValues() (ReduxValues, error) {
	if err := rdx.kvr.IndexRefresh(); err != nil {
		return rdx, err
	}

	rdxModTime, err := rdx.kvr.CurrentModTime(rdx.asset)
	if err != nil {
		return rdx, err
	}

	if rdx.connTime < rdxModTime {
		return ConnectRedux(rdx.dir, rdx.asset)
	}

	return rdx, nil
}

func (rdx *redux) Match(terms []string, scope map[string]bool, anyCase bool, contains bool) map[string]bool {
	if scope == nil {
		scope = make(map[string]bool)
		for _, k := range rdx.Keys() {
			scope[k] = true
		}
	}

	matches := make(map[string]bool)
	for _, term := range terms {
		if anyCase {
			term = strings.ToLower(term)
		}
		for key := range scope {
			if values, ok := rdx.GetAllValues(key); !ok {
				continue
			} else if anyValueMatchesTerm(term, values, anyCase, contains) {
				matches[key] = true
			}
		}
	}

	return matches
}

func anyValueMatchesTerm(term string, values []string, anyCase bool, contains bool) bool {
	for _, val := range values {
		if anyCase {
			val = strings.ToLower(val)
		}
		if contains {
			if strings.Contains(val, term) {
				return true
			}
		} else {
			if val == term {
				return true
			}
		}
	}
	return false
}
