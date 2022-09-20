package kvas

import (
	"fmt"
	"sort"
	"sync"
)

type reduxList struct {
	reductions map[string]ReduxValues
	fabric     *ReduxFabric
	mtx        *sync.Mutex
	modTime    int64
}

func ConnectReduxAssets(dir string, fabric *ReduxFabric, assets ...string) (ReduxAssets, error) {

	reductions := make(map[string]ReduxValues)
	var err error

	fabric = initFabric(fabric)

	for d := range fabric.Aggregates.DetailAll(assets...) {
		reductions[d], err = ConnectRedux(dir, d)
		if err != nil {
			return nil, err
		}

		if fabric.Transitives.IsTransitive(d) {
			td := fabric.Transitives.Transition(d)
			if _, ok := reductions[td]; !ok {
				reductions[td], err = ConnectRedux(dir, td)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return &reduxList{
		mtx:        &sync.Mutex{},
		fabric:     fabric,
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

func (rl *reduxList) AddVal(asset, key, val string) error {
	if !rl.Has(asset) {
		return fmt.Errorf("asset %s is not present in this list", asset)
	}
	return rl.reductions[asset].AddVal(key, val)
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

func (rl *reduxList) transitionValues(asset string, values ...string) []string {
	if rl.fabric == nil ||
		rl.fabric.Transitives == nil {
		return values
	}
	if rl.fabric.Transitives.IsTransitive(asset) {
		tk := rl.fabric.Transitives.Transition(asset)
		for i := 0; i < len(values); i++ {
			tv, ok := rl.reductions[tk].GetFirstVal(values[i])
			if !ok {
				continue
			}
			values[i] = rl.fabric.Transitives.Fmt(values[i], tv)
		}
	}
	return values
}

func (rl *reduxList) GetFirstVal(asset, key string) (string, bool) {
	if !rl.Has(asset) {
		return "", false
	}
	value, ok := rl.reductions[asset].GetFirstVal(key)
	tvs := rl.transitionValues(asset, value)
	if len(tvs) > 0 {
		value = tvs[0]
	}
	return value, ok
}

func (rl *reduxList) GetAllUnchangedValues(asset, key string) ([]string, bool) {
	if _, ok := rl.reductions[asset]; !ok {
		return nil, false
	}
	return rl.reductions[asset].GetAllValues(key)
}

func (rl *reduxList) GetAllValues(asset, key string) ([]string, bool) {
	values, ok := rl.GetAllUnchangedValues(asset, key)
	return rl.transitionValues(asset, values...), ok
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

// appendReverseReplacedTerms adds reversed transitioned (original) terms
// for a (transitioned) property.
// Example: pr-id is transitive with pr-name: pr-id: "1", pr-name: "property_one"
// would result in pr-id value displayed in "property_one (1)".
// Matching example: for a query {pr-id: {"property_one"}}, appendReverseReplacedTerms
// would transform that to {pr-id: {"property_one", "1"}} and objects that have
// pr-id:"1" would match.
// Note: reverse transitions automatically take atomicity into account.
func (rl *reduxList) appendReverseTransitions(asset string, terms []string, anyCase bool) []string {
	if rl.fabric.Transitives.IsTransitive(asset) {
		rp := rl.fabric.Transitives.Transition(asset)
		atomic := rl.fabric.Atomics.IsAtomic(rp)
		sourceTerms := rl.reductions[rp].Match(terms, nil, anyCase, !atomic)
		for t := range sourceTerms {
			terms = append(terms, t)
		}
	}
	return terms
}

// matchDetailed
func (rl *reduxList) matchDetailed(asset string, scope map[string]bool, terms []string, anyCase bool) map[string]bool {
	details := rl.fabric.Aggregates.Detail(asset)
	matches := make(map[string]bool, 0)
	for _, da := range details {
		terms = rl.appendReverseTransitions(da, terms, anyCase)
		atomic := rl.fabric.Atomics.IsAtomic(asset)
		results := rl.reductions[da].Match(terms, scope, anyCase, !atomic)
		for key := range results {
			if !matches[key] {
				matches[key] = true
			}
		}
	}
	return matches
}

func (rl *reduxList) Match(query map[string][]string, anyCase bool) map[string]bool {
	var matches map[string]bool
	for asset, terms := range query {
		if rl.fabric.Aggregates.IsAggregate(asset) {
			matches = rl.matchDetailed(asset, matches, terms, anyCase)
		} else {
			if _, ok := rl.reductions[asset]; !ok {
				continue
			}
			atomic := rl.fabric.Atomics.IsAtomic(asset)
			matches = rl.reductions[asset].Match(
				rl.appendReverseTransitions(asset, terms, anyCase),
				matches,
				anyCase,
				!atomic)
		}
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

type idPropertiesTitle struct {
	id         string
	properties []string
}

type sortableIdSet struct {
	ipt []idPropertiesTitle
}

func (is *sortableIdSet) Len() int {
	return len(is.ipt)
}

func (is *sortableIdSet) Swap(i, j int) {
	is.ipt[i], is.ipt[j] = is.ipt[j], is.ipt[i]
}

func (is *sortableIdSet) Less(i, j int) bool {
	for p, _ := range is.ipt[i].properties {
		if is.ipt[i].properties[p] < is.ipt[j].properties[p] {
			return true
		}
	}
	return false
}

func (rl *reduxList) Sort(ids []string, desc bool, sortBy ...string) ([]string, error) {
	if err := rl.IsSupported(sortBy...); err != nil {
		return nil, err
	}

	sis := &sortableIdSet{
		ipt: make([]idPropertiesTitle, 0, len(ids)),
	}

	for _, p := range sortBy {
		for _, id := range ids {
			ipt := idPropertiesTitle{id: id}
			v, _ := rl.GetFirstVal(p, id)
			ipt.properties = append(ipt.properties, v)
			sis.ipt = append(sis.ipt, ipt)
		}
	}

	var sortInterface sort.Interface = sis
	if desc {
		sortInterface = sort.Reverse(sortInterface)
	}

	sort.Sort(sortInterface)

	sorted := make([]string, 0, len(sis.ipt))
	for _, ipt := range sis.ipt {
		sorted = append(sorted, ipt.id)
	}

	return sorted, nil
}
