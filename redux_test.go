package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	testAsset   = "test_asset"
	detailAsset = "detail_asset"
)

func reduxCleanup(assets ...string) error {
	for _, asset := range assets {
		rdxPath := filepath.Join(os.TempDir(), asset+GobExt)
		if _, err := os.Stat(rdxPath); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if err := os.Remove(rdxPath); err != nil {
			return err
		}
	}

	return indexCleanup()
}

func mockRedux() *redux {
	return &redux{
		dir:   os.TempDir(),
		asset: testAsset,
		keyReductions: map[string][]string{
			"k1": {"v1", "v2", "v3"},
			"k2": {"v11", "v12", "v13"},
			"k3": {"v121", "v122", "v123"},
			"k4": {"v1231", "v1232", "v1233"},
		},
	}
}

func mockDetailsRedux() *redux {
	return &redux{
		dir:   os.TempDir(),
		asset: detailAsset,
		keyReductions: map[string][]string{
			"v1": {"d1", "d2"},
			"v2": {"d21", "d22"},
			"v3": {"d31", "d32"},
		},
	}
}

func TestReduxWriteConnect(t *testing.T) {
	wrdx := &redux{
		dir:           os.TempDir(),
		asset:         testAsset,
		keyReductions: map[string][]string{},
	}

	testo.Error(t, wrdx.write(), false)

	rdx, err := ConnectRedux(os.TempDir(), testAsset)
	testo.Error(t, err, false)
	testo.Nil(t, rdx, false)

	testo.Error(t, reduxCleanup(testAsset), false)
}

func TestReduxKeys(t *testing.T) {
	rdx := mockRedux()

	keys := rdx.Keys()
	testo.EqualValues(t, len(keys), len(rdx.keyReductions))
	for _, k := range keys {
		_, ok := rdx.keyReductions[k]
		testo.EqualValues(t, ok, true)
	}
}

func TestReduxHas(t *testing.T) {
	tests := []struct {
		key string
		exp bool
	}{
		{"", false},
		{"k1", true},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.EqualValues(t, rdx.Has(tt.key), tt.exp)
		})
	}
}

func TestReduxHasVal(t *testing.T) {
	tests := []struct {
		key, val string
		exp      bool
	}{
		{"", "", false},
		{"k1", "", false},
		{"k1", "v1", true},
		{"k1", "v11", false},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.key+tt.val, func(t *testing.T) {
			testo.EqualValues(t, rdx.HasVal(tt.key, tt.val), tt.exp)
		})
	}
}

func TestReduxAddVal(t *testing.T) {
	tests := []struct {
		key, val string
		exist    bool
	}{
		{"", "new", false},
		{"k1", "", false},
		{"k1", "v1", true},
		{"k1", "v11", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+tt.val, func(t *testing.T) {
			rdx := mockRedux()
			testo.EqualValues(t, rdx.HasVal(tt.key, tt.val), tt.exist)
			testo.Error(t, rdx.AddValues(tt.key, tt.val), false)
			testo.EqualValues(t, rdx.HasVal(tt.key, tt.val), true)
			testo.Error(t, reduxCleanup(testAsset), false)
		})
	}
}

func TestReduxReplaceValues(t *testing.T) {
	tests := []struct {
		key    string
		values []string
	}{
		{"", []string{"new"}},
		{"k1", nil},
		{"k1", []string{""}},
		{"k1", []string{"v1", "v2"}},
	}

	for _, tt := range tests {
		t.Run(tt.key+strings.Join(tt.values, ""), func(t *testing.T) {
			rdx := mockRedux()
			testo.Error(t, rdx.ReplaceValues(tt.key, tt.values...), false)
			testo.DeepEqual(t, rdx.keyReductions[tt.key], tt.values)
			testo.Error(t, reduxCleanup(testAsset), false)
		})
	}
}

func TestReduxBatchReplaceValues(t *testing.T) {
	tests := []struct {
		keyValues map[string][]string
	}{
		{map[string][]string{"": {"new"}}},
		{map[string][]string{"k1": nil}},
		{map[string][]string{"k1": {""}}},
		{map[string][]string{"k1": {"v1", "v2"}}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			rdx := mockRedux()
			testo.Error(t, rdx.BatchReplaceValues(tt.keyValues), false)
			for key, values := range tt.keyValues {
				testo.DeepEqual(t, rdx.keyReductions[key], values)
			}
			testo.Error(t, reduxCleanup(testAsset), false)
		})
	}
}

func TestReduxCutVal(t *testing.T) {
	tests := []struct {
		key, val string
		exist    bool
	}{
		{"", "new", false},
		{"k1", "", false},
		{"k1", "v1", true},
	}

	for _, tt := range tests {
		t.Run(tt.key+tt.val, func(t *testing.T) {
			rdx := mockRedux()
			testo.EqualValues(t, rdx.HasVal(tt.key, tt.val), tt.exist)
			testo.Error(t, rdx.CutVal(tt.key, tt.val), false)
			testo.EqualValues(t, rdx.HasVal(tt.key, tt.val), false)
			testo.Error(t, reduxCleanup(testAsset), false)
		})
	}
}

func TestReduxGetAllValues(t *testing.T) {
	tests := []struct {
		key string
		ok  bool
	}{
		{"", false},
		{"k1", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			rdx := mockRedux()
			values, ok := rdx.GetAllValues(tt.key)
			if ok {
				testo.DeepEqual(t, rdx.keyReductions[tt.key], values)
			}
			testo.EqualValues(t, ok, tt.ok)
			testo.Error(t, reduxCleanup(testAsset), false)
		})
	}
}

func TestReduxGetFirstVal(t *testing.T) {
	tests := []struct {
		key string
		ok  bool
	}{
		{"", false},
		{"k1", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			rdx := mockRedux()
			fv, ok := rdx.GetFirstVal(tt.key)
			if ok {
				testo.DeepEqual(t, rdx.keyReductions[tt.key][0], fv)
			}
			testo.EqualValues(t, ok, tt.ok)
			testo.Error(t, reduxCleanup(testAsset), false)
		})
	}
}

func TestAnyValueMatchesTerm(t *testing.T) {
	tests := []struct {
		term     string
		values   []string
		anyCase  bool
		contains bool
		ok       bool
	}{
		{"test", []string{"nomatch", "prefixTEST"}, false, false, false},
		{"prefixTEST", []string{"nomatch", "prefixTEST"}, false, false, true},
		{"test", []string{"nomatch", "prefixTEST"}, true, false, false},
		{"prefixtest", []string{"nomatch", "prefixTEST"}, true, false, true},
		{"test", []string{"nomatch", "prefixTEST"}, false, true, false},
		{"test", []string{"nomatch", "prefixTEST"}, true, true, true},
		{"test", []string{"nomatch"}, true, true, false},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			ok := anyValueMatchesTerm(tt.term, tt.values, tt.anyCase, tt.contains)
			testo.EqualValues(t, ok, tt.ok)
		})
	}
}

func TestReduxMatch(t *testing.T) {
	tests := []struct {
		terms    []string
		scope    map[string]bool
		anyCase  bool
		contains bool
		matches  []string
	}{
		{[]string{"11"}, nil, false, false, []string{}},
		{[]string{"11"}, nil, false, true, []string{"k2"}},
		{[]string{"11"}, nil, true, true, []string{"k2"}},
		{[]string{"11"}, map[string]bool{"k1": true, "k3": true}, true, true, []string{}},
		{[]string{"V12"}, nil, false, false, []string{}},
		{[]string{"V12"}, nil, true, false, []string{"k2"}},
		{[]string{"V12"}, nil, false, true, []string{}},
		{[]string{"V12"}, nil, true, true, []string{"k2", "k3", "k4"}},
		{[]string{"V12"}, map[string]bool{"k4": true, "k5": true}, true, true, []string{"k4"}},
	}

	rdx := mockRedux()

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {

			matches := rdx.Match(tt.terms, tt.scope, tt.anyCase, tt.contains)
			testo.EqualValues(t, len(matches), len(tt.matches))
			for _, m := range tt.matches {
				_, ok := matches[m]
				testo.EqualValues(t, ok, true)
			}
		})
	}
}
