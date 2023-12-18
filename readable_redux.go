package kvas

import (
	"github.com/boggydigital/wits"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"io"
	"sort"
	"strings"
)

func (rdx *Redux) MustHave(assets ...string) error {
	for _, asset := range assets {
		if !rdx.Has(asset) {
			return UnknownReduxAsset(asset)
		}
	}
	return nil
}

func (rdx *Redux) Keys(asset string) []string {
	return maps.Keys(rdx.assetKeyValues[asset])
}

func (rdx *Redux) Has(asset string) bool {
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

func (rdx *Redux) HasVal(asset, key, val string) bool {
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

func (rdx *Redux) GetFirstVal(asset, key string) (string, bool) {
	if values, ok := rdx.GetAllValues(asset, key); ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

func (rdx *Redux) GetAllValues(asset, key string) ([]string, bool) {
	if !rdx.Has(asset) {
		return nil, false
	}
	if rdx.assetKeyValues[asset] == nil {
		return nil, false
	}

	val, ok := rdx.assetKeyValues[asset][key]
	return val, ok
}

func (rdx *Redux) assetModTime(asset string) (int64, error) {
	return rdx.kvr.CurrentModTime(asset)
}

func (rdx *Redux) ModTime() (int64, error) {
	rdx.mtx.Lock()
	defer rdx.mtx.Unlock()

	var mt int64 = 0
	for asset := range rdx.assetKeyValues {
		if amt, err := rdx.assetModTime(asset); err != nil {
			return -1, err
		} else {
			if mt < amt {
				mt = amt
			}
		}
	}
	return mt, nil
}

func (rdx *Redux) Refresh() error {
	if err := rdx.kvr.IndexRefresh(); err != nil {
		return err
	}

	modTime, err := rdx.ModTime()
	if err != nil {
		return err
	}

	if rdx.modTime < modTime {
		if rdx, err = connectRedux(rdx.dir, maps.Keys(rdx.assetKeyValues)...); err != nil {
			return err
		}
	}

	return nil
}

func (rdx *Redux) MatchAsset(asset string, terms []string, scope []string, anyCase bool, contains bool) []string {
	if scope == nil {
		scope = rdx.Keys(asset)
	}

	matches := make(map[string]interface{})
	for _, term := range terms {
		if anyCase {
			term = strings.ToLower(term)
		}
		for _, key := range scope {
			if values, ok := rdx.GetAllValues(asset, key); !ok {
				continue
			} else if anyValueMatchesTerm(term, values, anyCase, contains) {
				matches[key] = nil
			}
		}
	}

	return maps.Keys(matches)
}

func (rdx *Redux) Match(query map[string][]string, anyCase, contains bool) []string {
	var matches []string
	for asset, terms := range query {
		if !rdx.Has(asset) {
			continue
		}
		matches = rdx.MatchAsset(asset, terms, matches, anyCase, contains)
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

func (rdx *Redux) Sort(ids []string, desc bool, sortBy ...string) ([]string, error) {
	if err := rdx.MustHave(sortBy...); err != nil {
		return nil, err
	}

	sis := &sortableIdSet{
		properties: sortBy,
		ipv:        make([]idValues, 0, len(ids)),
	}

	for _, id := range ids {
		iv := idValues{id: id}
		for _, p := range sortBy {
			v, _ := rdx.GetFirstVal(p, id)
			iv.values = append(iv.values, v)
		}
		sis.ipv = append(sis.ipv, iv)
	}

	var sortInterface sort.Interface = sis
	if desc {
		sortInterface = sort.Reverse(sortInterface)
	}

	sort.Sort(sortInterface)

	sorted := make([]string, 0, len(sis.ipv))
	for _, iv := range sis.ipv {
		sorted = append(sorted, iv.id)
	}

	return sorted, nil
}

func (rdx *Redux) Export(w io.Writer, ids ...string) error {

	skv := make(wits.SectionKeyValues)

	for _, id := range ids {
		skv[id] = make(wits.KeyValues)
		for p := range rdx.assetKeyValues {
			if vals, ok := rdx.GetAllValues(p, id); ok {
				skv[id][p] = vals
			}
		}
	}

	return skv.Write(w)
}
