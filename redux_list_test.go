package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"sort"
	"strconv"
	"testing"
)

var mockAssets = []string{"a1", "a2"}

func mockReduxList() *reduxList {
	return &reduxList{
		assets: mockAssets,
		reductions: map[string]ReduxValues{
			"a1": mockRedux(),
			"a2": mockDetailsRedux(),
		},
		fabric: initFabric(nil),
	}
}

func reduxListCleanup() error {
	return reduxCleanup(append(mockAssets, testAsset, detailAsset)...)
}

func TestConnectReduxAssets(t *testing.T) {
	rxl, err := ConnectReduxAssets(os.TempDir(), nil, mockAssets...)
	testo.Error(t, err, false)
	testo.Nil(t, rxl, false)
	testo.Error(t, reduxListCleanup(), false)
}

func TestReduxListKeys(t *testing.T) {
	tests := []struct {
		asset string
		exp   []string
	}{
		{"", nil},
		{"a1", []string{"k1", "k2", "k3", "k4"}},
	}

	rxl := mockReduxList()
	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			keys := rxl.Keys(tt.asset)
			sort.Strings(keys)
			sort.Strings(tt.exp)
			testo.DeepEqual(t, keys, tt.exp)
		})
	}
}

func TestReduxListHas(t *testing.T) {
	tests := []struct {
		asset string
		exp   bool
	}{
		{"", false},
		{"a1", true},
	}

	rxl := mockReduxList()
	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			testo.EqualValues(t, rxl.Has(tt.asset), tt.exp)
		})
	}
}

func TestReduxListHasKey(t *testing.T) {
	tests := []struct {
		asset string
		key   string
		exp   bool
	}{
		{"", "", false},
		{"a1", "", false},
		{"a1", "k1", true},
	}

	rxl := mockReduxList()
	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			testo.EqualValues(t, rxl.HasKey(tt.asset, tt.key), tt.exp)
		})
	}
}

func TestReduxListHasVal(t *testing.T) {
	tests := []struct {
		asset string
		key   string
		val   string
		exp   bool
	}{
		{"", "", "", false},
		{"a1", "", "", false},
		{"a1", "k1", "", false},
		{"a1", "k1", "v1", true},
	}

	rxl := mockReduxList()
	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, tt.val), tt.exp)
		})
	}
}

func TestReduxListAddVal(t *testing.T) {
	tests := []struct {
		asset, key, val string
		has             bool
		expErr          bool
	}{
		{"", "k1", "v1", false, true},

		{"a1", "k1", "v1", true, false},
		{"a1", "k1", "v10", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.val, func(t *testing.T) {
			rxl := mockReduxList()
			testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, tt.val), tt.has)
			testo.Error(t, rxl.AddVal(tt.asset, tt.key, tt.val), tt.expErr)
			if !tt.expErr {
				testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, tt.val), true)
			}
			testo.Error(t, reduxListCleanup(), false)
		})
	}
}

func TestReduxListReplaceValues(t *testing.T) {
	tests := []struct {
		asset, key string
		values     []string
		expErr     bool
	}{
		{"", "", nil, true},
		{"a1", "k1", nil, false},
		{"a1", "k1", []string{"1", "2", "3"}, false},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			for _, nv := range tt.values {
				testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, nv), false)
			}
			testo.Error(t, rxl.ReplaceValues(tt.asset, tt.key, tt.values...), tt.expErr)
			if !tt.expErr {
				for _, nv := range tt.values {
					testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, nv), true)
				}
			}
			testo.Error(t, reduxListCleanup(), false)
		})
	}
}

func TestReduxListBatchReplaceValues(t *testing.T) {
	tests := []struct {
		asset     string
		keyValues map[string][]string
		expErr    bool
	}{
		{"", nil, true},
		{"a1", nil, false},
		{"a1", map[string][]string{"k1": nil, "k2": nil, "k3": nil}, false},
		{"a1", map[string][]string{"k1": {"v10"}, "k2": {"v20"}, "k3": {"v30"}}, false},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			for key, values := range tt.keyValues {
				for _, v := range values {
					testo.EqualValues(t, rxl.HasVal(tt.asset, key, v), false)
				}
			}
			testo.Error(t, rxl.BatchReplaceValues(tt.asset, tt.keyValues), tt.expErr)
			if !tt.expErr {
				for key, values := range tt.keyValues {
					for _, v := range values {
						testo.EqualValues(t, rxl.HasVal(tt.asset, key, v), true)
					}
				}
			}
			testo.Error(t, reduxListCleanup(), false)
		})
	}
}

func TestReduxListCutVal(t *testing.T) {
	tests := []struct {
		asset, key, val string
		has             bool
		expErr          bool
	}{
		{"", "", "v1", false, true},
		{"a1", "k1", "v1", true, false},
		{"a1", "k1", "v10", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.val, func(t *testing.T) {
			rxl := mockReduxList()
			testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, tt.val), tt.has)
			testo.Error(t, rxl.CutVal(tt.asset, tt.key, tt.val), tt.expErr)
			if !tt.expErr {
				testo.EqualValues(t, rxl.HasVal(tt.asset, tt.key, tt.val), false)
			}
			testo.Error(t, reduxListCleanup(), false)
		})
	}
}

func TestReduxListTransitionValues(t *testing.T) {
	tests := []struct {
		asset       string
		transitives map[string]string
		values      []string
		exp         []string
	}{
		{"a1", nil, nil, nil},
		{
			"a1",
			map[string]string{"a2": "a1"},
			[]string{"v1", "v2"},
			[]string{"v1", "v2"},
		},
		{
			"a1",
			map[string]string{"a1": "a1"},
			[]string{"v1", "v2"},
			[]string{"v1", "v2"},
		},
		{
			"a1",
			map[string]string{"a1": "a2"},
			[]string{"v1", "v2"},
			[]string{"d1 (v1)", "d21 (v2)"},
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Transitives = tt.transitives
			tv := rxl.transitionValues(tt.asset, tt.values...)
			testo.DeepEqual(t, tv, tt.exp)
		})
	}
}

func TestReduxListGetFirstVal(t *testing.T) {
	tests := []struct {
		transitives map[string]string
		asset, key  string
		expVal      string
		expOk       bool
	}{
		{nil, "", "", "", false},
		{
			map[string]string{"": ""},
			"a1", "k1",
			"v1",
			true,
		},
		{
			map[string]string{"a1": "a2"},
			"a1", "k1",
			"d1 (v1)",
			true,
		},
		{
			map[string]string{"a1": "a2"},
			"a1", "unknown",
			"",
			false,
		},
		{
			map[string]string{"a1": "a2"},
			"unknown", "k1",
			"",
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Transitives = tt.transitives
			fv, ok := rxl.GetFirstVal(tt.asset, tt.key)
			testo.EqualValues(t, fv, tt.expVal)
			testo.EqualValues(t, ok, tt.expOk)
		})
	}
}

func TestReduxListGetAllUnchangedValues(t *testing.T) {
	tests := []struct {
		transitives map[string]string
		asset, key  string
		expValues   []string
		expOk       bool
	}{
		{nil, "", "", nil, false},
		{
			map[string]string{"": ""},
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			map[string]string{"a1": "a2"},
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			map[string]string{"a1": "a2"},
			"a1", "unknown",
			nil,
			false,
		},
		{
			map[string]string{"a1": "a2"},
			"unknown", "k1",
			nil,
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Transitives = tt.transitives
			fvs, ok := rxl.GetAllUnchangedValues(tt.asset, tt.key)
			testo.DeepEqual(t, fvs, tt.expValues)
			testo.EqualValues(t, ok, tt.expOk)
		})
	}
}

func TestReduxListGetAllValues(t *testing.T) {
	tests := []struct {
		transitives map[string]string
		asset, key  string
		expValues   []string
		expOk       bool
	}{
		{nil, "", "", nil, false},
		{
			map[string]string{"": ""},
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			map[string]string{"a1": "a2"},
			"a1", "k1",
			[]string{"d1 (v1)", "d21 (v2)", "d31 (v3)"},
			true,
		},
		{
			map[string]string{"a1": "a2"},
			"a1", "unknown",
			nil,
			false,
		},
		{
			map[string]string{"a1": "a2"},
			"unknown", "k1",
			nil,
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Transitives = tt.transitives
			fvs, ok := rxl.GetAllValues(tt.asset, tt.key)
			testo.DeepEqual(t, fvs, tt.expValues)
			testo.EqualValues(t, ok, tt.expOk)
		})
	}
}

func TestReduxListAppendReverseTransitions(t *testing.T) {
	tests := []struct {
		transitives map[string]string
		atomics     map[string]bool
		asset       string
		terms       []string
		anyCase     bool
		exp         []string
	}{
		{
			nil,
			nil,
			"a1",
			[]string{"d1", "d21"},
			true,
			[]string{"d1", "d21"},
		},
		{
			map[string]string{"a1": "a2"},
			nil,
			"a1",
			[]string{"d1", "d21"},
			true,
			[]string{"d1", "d21", "v1", "v2"},
		},
		{
			map[string]string{"a1": "a2"},
			nil,
			"a1",
			[]string{"d1", "x21"},
			true,
			[]string{"d1", "x21", "v1"},
		},
		{
			map[string]string{"a1": "a2"},
			nil,
			"a1",
			[]string{"d1", "d21"},
			false,
			[]string{"d1", "d21", "v1", "v2"},
		},
		{
			map[string]string{"a1": "a2"},
			nil,
			"a1",
			[]string{"d1", "D21"},
			false,
			[]string{"d1", "D21", "v1"},
		},
		{
			map[string]string{"a1": "a2"},
			nil,
			"a1",
			[]string{"d1", "d2"},
			false,
			[]string{"d1", "d2", "v1", "v2"},
		},
		{
			map[string]string{"a1": "a2"},
			map[string]bool{"a2": true},
			"a1",
			[]string{"d1", "d2"},
			false,
			[]string{"d1", "d2", "v1"},
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Transitives = tt.transitives
			rxl.fabric.Atomics = tt.atomics
			rtt := rxl.appendReverseTransitions(tt.asset, tt.terms, tt.anyCase)
			sort.Strings(rtt)
			sort.Strings(tt.exp)
			testo.DeepEqual(t, rtt, tt.exp)
		})
	}
}

func TestReduxListMatchDetailed(t *testing.T) {
	tests := []struct {
		aggregates map[string][]string
		atomics    map[string]bool
		asset      string
		scope      map[string]bool
		terms      []string
		anyCase    bool
		exp        map[string]bool
	}{
		{
			nil,
			nil,
			"a1",
			nil,
			[]string{"v1"},
			true,
			map[string]bool{},
		},
		{
			map[string][]string{"a1": {"a2"}},
			nil,
			"a1",
			nil,
			[]string{"v1"},
			true,
			map[string]bool{},
		},
		{
			map[string][]string{"a1": {"a2"}},
			nil,
			"a1",
			nil,
			[]string{"d1"},
			true,
			map[string]bool{"v1": true},
		},
		{
			map[string][]string{"a1": {"a1"}},
			nil,
			"a1",
			nil,
			[]string{"v1"},
			true,
			map[string]bool{"k1": true, "k2": true, "k3": true, "k4": true},
		},
		{
			map[string][]string{"a1": {"a1"}},
			map[string]bool{"a1": true},
			"a1",
			nil,
			[]string{"v1"},
			true,
			map[string]bool{"k1": true},
		},
		{
			map[string][]string{"a1": {"a1"}},
			map[string]bool{"a1": true},
			"a1",
			nil,
			[]string{"V1"},
			false,
			map[string]bool{},
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Aggregates = tt.aggregates
			rxl.fabric.Atomics = tt.atomics
			dm := rxl.matchDetailed(tt.asset, tt.scope, tt.terms, tt.anyCase)
			testo.EqualValues(t, len(dm), len(tt.exp))
			for dmt := range dm {
				_, ok := tt.exp[dmt]
				testo.EqualValues(t, ok, true)
			}
		})
	}
}

func TestReduxListMatch(t *testing.T) {
	tests := []struct {
		aggregates  map[string][]string
		transitives map[string]string
		atomics     map[string]bool
		query       map[string][]string
		anyCase     bool
		exp         map[string]bool
	}{
		{
			nil,
			nil,
			nil,
			nil,
			true,
			map[string]bool{},
		},
		{
			nil,
			nil,
			nil,
			map[string][]string{
				"a1": {"v2", "v123"},
			},
			true,
			map[string]bool{"k1": true, "k3": true, "k4": true},
		},
		{
			nil,
			nil,
			nil,
			map[string][]string{
				"a1": {"V2", "v123"},
			},
			false,
			map[string]bool{"k3": true, "k4": true},
		},
		{
			nil,
			map[string]string{"a1": "a2"},
			nil,
			map[string][]string{
				"a1": {"d2"},
			},
			false,
			map[string]bool{"k1": true, "k2": true, "k3": true, "k4": true},
		},
		{
			nil,
			map[string]string{"a1": "a2"},
			map[string]bool{"a1": true},
			map[string][]string{
				"a1": {"d2"},
			},
			false,
			map[string]bool{"k1": true},
		},
		{
			map[string][]string{"a1": {"a2"}},
			nil,
			nil,
			map[string][]string{
				"a1": {"d1"},
			},
			false,
			map[string]bool{"v1": true},
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Aggregates = tt.aggregates
			rxl.fabric.Transitives = tt.transitives
			rxl.fabric.Atomics = tt.atomics
			ms := rxl.Match(tt.query, tt.anyCase)
			testo.EqualValues(t, len(ms), len(tt.exp))
			for mt := range ms {
				_, ok := tt.exp[mt]
				testo.EqualValues(t, ok, true)
			}
		})
	}

}

func TestReduxListIsSupported(t *testing.T) {
	tests := []struct {
		asset  string
		expErr bool
	}{
		{"", true},
		{"a1", false},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			rxl := mockReduxList()
			testo.Error(t, rxl.IsSupported(tt.asset), tt.expErr)
		})
	}
}
