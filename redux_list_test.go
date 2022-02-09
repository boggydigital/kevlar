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

func TestConnectReduxList(t *testing.T) {
	rxl, err := ConnectReduxList(os.TempDir(), nil, mockAssets...)
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
		transitives map[string]string
		values      []string
		exp         []string
	}{
		{nil, nil, nil},
		{
			map[string]string{"a2": "a1"},
			[]string{"v1", "v2"},
			[]string{"v1", "v2"},
		},
		{
			map[string]string{"a1": "a1"},
			[]string{"v1", "v2"},
			[]string{"v1", "v2"},
		},
		{
			map[string]string{"a1": "a2"},
			[]string{"v1", "v2"},
			[]string{"d1 (v1)", "d21 (v2)"},
		},
	}

	asset := "a1"
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			rxl.fabric.Transitives = tt.transitives
			tv := rxl.transitionValues(asset, tt.values...)
			testo.DeepEqual(t, tv, tt.exp)
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
