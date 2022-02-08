package kvas

import (
	"bytes"
	"encoding/gob"
	"strings"
)

type redux struct {
	dir           string
	asset         string
	keyReductions map[string][]string
}

func ConnectRedux(dir, asset string) (ReduxValues, error) {
	rdx, err := ConnectLocal(dir, GobExt)
	if err != nil {
		return nil, err
	}

	arc, err := rdx.Get(asset)
	defer arc.Close()
	if err != nil {
		return nil, err
	}

	var keyReductions map[string][]string
	if err := gob.NewDecoder(arc).Decode(&keyReductions); err != nil {
		return nil, err
	}

	if keyReductions == nil {
		keyReductions = make(map[string][]string, 0)
	}

	return &redux{
		dir:           dir,
		asset:         asset,
		keyReductions: keyReductions,
	}, nil
}

func (rdx *redux) Keys() []string {
	keys := make([]string, 0, len(rdx.keyReductions))
	for k := range rdx.keyReductions {
		keys = append(keys, k)
	}
	return keys
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

func (rdx *redux) AddVal(key string, val string) error {
	if rdx.HasVal(key, val) {
		return nil
	}
	rdx.keyReductions[key] = append(rdx.keyReductions[key], val)
	return rdx.write()
}

func (rdx *redux) ReplaceValues(key string, values ...string) error {
	rdx.keyReductions[key] = values
	return rdx.write()
}

func (rdx *redux) BatchReplaceValues(keysValues map[string][]string) error {
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

func (rdx *redux) GetAllValues(key string) ([]string, bool) {
	if rdx == nil || rdx.keyReductions == nil {
		return nil, false
	}
	val, ok := rdx.keyReductions[key]
	return val, ok
}

func (rdx *redux) GetFirstVal(key string) (string, bool) {
	values, ok := rdx.GetAllValues(key)
	if ok && len(values) > 0 {
		return values[0], ok
	}
	return "", false
}

func (rdx *redux) Match(terms []string, scope []string, anyCase bool, contains bool) map[string]bool {
	if scope == nil {
		scope = rdx.Keys()
	}

	matches := make(map[string]bool)
	for _, term := range terms {
		if anyCase {
			term = strings.ToLower(term)
		}
		for _, key := range scope {
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
