package kvas

import (
	"fmt"
	"github.com/boggydigital/wits"
	"io"
	"sort"
	"sync"
)

type reduxList struct {
	reductions map[string]ReduxValues
	mtx        *sync.Mutex
	modTime    int64
}

func ConnectReduxAssets(dir string, assets ...string) (ReduxAssets, error) {

	reductions := make(map[string]ReduxValues)
	var err error

	for _, d := range assets {
		reductions[d], err = ConnectRedux(dir, d)
		if err != nil {
			return nil, err
		}
	}

	return &reduxList{
		mtx:        &sync.Mutex{},
		reductions: reductions,
		modTime:    -1,
	}, nil
}

func (rl *reduxList) Keys(asset string) []string {
	if _, ok := rl.reductions[asset]; !ok {
		return nil
	}
	return rl.reductions[asset].Keys()
}

func (rl *reduxList) Has(asset string) bool {
	_, ok := rl.reductions[asset]
	return ok
}

func (rl *reduxList) HasKey(asset, key string) bool {
	if !rl.Has(asset) {
		return false
	}
	return rl.reductions[asset].Has(key)
}

func (rl *reduxList) HasVal(asset, key, val string) bool {
	if !rl.Has(asset) {
		return false
	}
	return rl.reductions[asset].HasVal(key, val)
}

func (rl *reduxList) AddValues(asset, key string, values ...string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].AddValues(key, values...)
}

func (rl *reduxList) BatchAddValues(asset string, keyValues map[string][]string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].BatchAddValues(keyValues)
}

func (rl *reduxList) ReplaceValues(asset, key string, values ...string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].ReplaceValues(key, values...)
}

func (rl *reduxList) BatchReplaceValues(asset string, keyValues map[string][]string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].BatchReplaceValues(keyValues)
}

func (rl *reduxList) CutVal(asset, key, val string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].CutVal(key, val)
}

func (rl *reduxList) BatchCutValues(asset string, keyValues map[string][]string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].BatchCutValues(keyValues)
}

func (rl *reduxList) GetFirstVal(asset, key string) (string, bool) {
	if !rl.Has(asset) {
		return "", false
	}
	return rl.reductions[asset].GetFirstVal(key)
}

func (rl *reduxList) GetAllValues(asset, key string) ([]string, bool) {
	if _, ok := rl.reductions[asset]; !ok {
		return nil, false
	}
	return rl.reductions[asset].GetAllValues(key)
}

func (rl *reduxList) RefreshReduxAssets() (ReduxAssets, error) {

	modTime, err := rl.ReduxAssetsModTime()
	if err != nil {
		return rl, err
	}

	if rl.modTime >= modTime {
		return rl, nil
	}

	for asset := range rl.reductions {
		if rl.reductions[asset], err = rl.reductions[asset].RefreshReduxValues(); err != nil {
			return rl, err
		}
	}

	rl.modTime = modTime

	return rl, nil
}

func (rl *reduxList) ReduxAssetsModTime() (int64, error) {
	rl.mtx.Lock()
	defer rl.mtx.Unlock()

	mt := int64(0)
	for _, rdx := range rl.reductions {
		if rmt, err := rdx.ReduxModTime(); err != nil {
			return mt, err
		} else {
			if mt < rmt {
				mt = rmt
			}
		}
	}
	return mt, nil
}

func (rl *reduxList) Match(query map[string][]string, anyCase, contains bool) map[string]bool {
	var matches map[string]bool
	for asset, terms := range query {
		if _, ok := rl.reductions[asset]; !ok {
			continue
		}
		matches = rl.reductions[asset].Match(
			terms,
			matches,
			anyCase,
			contains)
	}
	return matches
}

func (rl *reduxList) IsSupported(assets ...string) error {
	for _, a := range assets {
		if _, ok := rl.reductions[a]; !ok {
			return fmt.Errorf("unsupported asset %s", a)
		}
	}

	return nil
}

type idValues struct {
	id     string
	values []string
}

type sortableIdSet struct {
	properties []string
	ipv        []idValues
}

func (is *sortableIdSet) Len() int {
	return len(is.ipv)
}

func (is *sortableIdSet) Swap(i, j int) {
	is.ipv[i], is.ipv[j] = is.ipv[j], is.ipv[i]
}

func (is *sortableIdSet) Less(i, j int) bool {
	for p, _ := range is.properties {
		if is.ipv[i].values[p] == is.ipv[j].values[p] {
			continue
		}
		return is.ipv[i].values[p] < is.ipv[j].values[p]
	}
	return false
}

func (rl *reduxList) Sort(ids []string, desc bool, sortBy ...string) ([]string, error) {
	if err := rl.IsSupported(sortBy...); err != nil {
		return nil, err
	}

	sis := &sortableIdSet{
		properties: sortBy,
		ipv:        make([]idValues, 0, len(ids)),
	}

	for _, id := range ids {
		iv := idValues{id: id}
		for _, p := range sortBy {
			v, _ := rl.GetFirstVal(p, id)
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

func (rl *reduxList) Export(w io.Writer, ids ...string) error {

	skv := make(wits.SectionKeyValues)

	for _, id := range ids {
		skv[id] = make(wits.KeyValues)
		for p := range rl.reductions {
			if vals, ok := rl.GetAllValues(p, id); ok {
				skv[id][p] = vals
			}
		}
	}

	return skv.Write(w)
}
