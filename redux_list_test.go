package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"sort"
	"strconv"
	"testing"
)

func mockReduxList() *reduxList {
	return &reduxList{
		reductions: map[string]ReduxValues{
			"a1": mockRedux(),
			"a2": mockDetailsRedux(),
		},
	}
}

func mockAssets() []string {
	assets := []string{}
	rxa := mockReduxList()
	for a := range rxa.reductions {
		assets = append(assets, a)
	}
	return assets
}

func reduxListCleanup() error {
	return reduxCleanup(append(mockAssets(), testAsset, detailAsset)...)
}

func TestConnectReduxAssets(t *testing.T) {
	tests := []struct {
		assets        []string
		connectionErr bool
	}{
		{mockAssets(), false},
		{[]string{"a2"}, false},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl, err := ConnectReduxAssets(os.TempDir(), tt.assets...)
			testo.Error(t, err, tt.connectionErr)
			testo.Nil(t, rxl, false)
			for _, a := range tt.assets {
				testo.Error(t, rxl.IsSupported(a), false)
			}
			testo.Error(t, reduxListCleanup(), false)
		})
	}

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

func TestReduxListGetFirstVal(t *testing.T) {
	tests := []struct {
		asset, key string
		expVal     string
		expOk      bool
	}{
		{"", "", "", false},
		{
			"a1", "k1",
			"v1",
			true,
		},
		{
			"a1", "k1",
			"v1",
			true,
		},
		{
			"a1", "unknown",
			"",
			false,
		},
		{
			"unknown", "k1",
			"",
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			fv, ok := rxl.GetFirstVal(tt.asset, tt.key)
			testo.EqualValues(t, fv, tt.expVal)
			testo.EqualValues(t, ok, tt.expOk)
		})
	}
}

func TestReduxListGetAllUnchangedValues(t *testing.T) {
	tests := []struct {
		asset, key string
		expValues  []string
		expOk      bool
	}{
		{"", "", nil, false},
		{
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			"a1", "unknown",
			nil,
			false,
		},
		{
			"unknown", "k1",
			nil,
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			fvs, ok := rxl.GetAllValues(tt.asset, tt.key)
			testo.DeepEqual(t, fvs, tt.expValues)
			testo.EqualValues(t, ok, tt.expOk)
		})
	}
}

func TestReduxListGetAllValues(t *testing.T) {
	tests := []struct {
		asset, key string
		expValues  []string
		expOk      bool
	}{
		{"", "", nil, false},
		{
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			"a1", "k1",
			[]string{"v1", "v2", "v3"},
			true,
		},
		{
			"a1", "unknown",
			nil,
			false,
		},
		{
			"unknown", "k1",
			nil,
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			fvs, ok := rxl.GetAllValues(tt.asset, tt.key)
			testo.DeepEqual(t, fvs, tt.expValues)
			testo.EqualValues(t, ok, tt.expOk)
		})
	}
}

func TestReduxListMatch(t *testing.T) {
	tests := []struct {
		query    map[string][]string
		anyCase  bool
		contains bool
		exp      map[string]bool
	}{
		{
			nil,
			true,
			true,
			map[string]bool{},
		},
		{
			map[string][]string{
				"a1": {"v2", "v123"},
			},
			true,
			true,
			map[string]bool{"k1": true, "k3": true, "k4": true},
		},
		{
			map[string][]string{
				"a1": {"V2", "v123"},
			},
			false,
			true,
			map[string]bool{"k3": true, "k4": true},
		},
		{
			map[string][]string{
				"a1": {"v2"},
			},
			false,
			true,
			map[string]bool{"k1": true},
		},
		{
			map[string][]string{
				"a1": {"d2"},
			},
			false,
			false,
			map[string]bool{},
		},
		{
			map[string][]string{
				"a1": {"d1"},
			},
			false,
			true,
			map[string]bool{},
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rxl := mockReduxList()
			ms := rxl.Match(tt.query, tt.anyCase, tt.contains)
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
