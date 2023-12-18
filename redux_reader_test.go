package kvas

import (
	"github.com/boggydigital/testo"
	"strconv"
	"testing"
)

func TestRedux_MustHave(t *testing.T) {
	tests := []struct {
		assets []string
		errExp bool
	}{
		{[]string{}, false},
		{[]string{""}, true},
		{[]string{"a1"}, false},
		{[]string{"a1", "1a"}, true},
	}

	rdx := mockRedux()
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			err := rdx.MustHave(tt.assets...)
			testo.Error(t, err, tt.errExp)
		})
	}
}

func TestReduxKeys(t *testing.T) {
	rdx := mockRedux()
	for asset := range rdx.assetKeyValues {
		keys := rdx.Keys(asset)
		testo.EqualValues(t, len(keys), len(rdx.assetKeyValues[asset]))
		for _, k := range keys {
			_, ok := rdx.assetKeyValues[asset][k]
			testo.EqualValues(t, ok, true)
		}
	}
}

func TestReduxHas(t *testing.T) {
	tests := []struct {
		asset string
		exp   bool
	}{
		{"", false},
		{"a1", true},
		{"1a", false},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.asset, func(t *testing.T) {
			testo.EqualValues(t, rdx.Has(tt.asset), tt.exp)
		})
	}
}

func TestReduxHasKey(t *testing.T) {
	tests := []struct {
		asset, key string
		exp        bool
	}{
		{"", "", false},
		{"a1", "", false},
		{"a1", "k1", true},
		{"a1", "1k", false},
		{"1a", "k1", false},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.asset+tt.key, func(t *testing.T) {
			testo.EqualValues(t, rdx.HasKey(tt.asset, tt.key), tt.exp)
		})
	}
}

func TestReduxHasValue(t *testing.T) {
	tests := []struct {
		asset, key, value string
		exp               bool
	}{
		{"", "", "", false},
		{"a1", "", "", false},
		{"a1", "k1", "", false},
		{"a1", "k1", "v11", true},
		{"1a", "k1", "v11", false},
		{"a1", "1k", "v11", false},
		{"a1", "k1", "11v", false},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.asset+tt.key+tt.value, func(t *testing.T) {
			testo.EqualValues(t, rdx.HasValue(tt.asset, tt.key, tt.value), tt.exp)
		})
	}
}

func TestReduxGetAllValues(t *testing.T) {
	tests := []struct {
		asset, key string
		ok         bool
	}{
		{"", "", false},
		{"a1", "k1", true},
		{"1a", "k1", false},
		{"a1", "1k", false},
	}

	for _, tt := range tests {
		t.Run(tt.asset+tt.key, func(t *testing.T) {
			rdx := mockRedux()
			values, ok := rdx.GetAllValues(tt.asset, tt.key)
			if ok {
				testo.DeepEqual(t, rdx.assetKeyValues[tt.asset][tt.key], values)
			}
			testo.EqualValues(t, ok, tt.ok)
		})
	}
}

func TestReduxGetFirstVal(t *testing.T) {
	tests := []struct {
		asset, key string
		ok         bool
	}{
		{"", "", false},
		{"a1", "k1", true},
		{"1a", "k1", false},
		{"a1", "1k", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			rdx := mockRedux()
			fv, ok := rdx.GetFirstVal(tt.asset, tt.key)
			if ok {
				testo.DeepEqual(t, rdx.assetKeyValues[tt.asset][tt.key][0], fv)
			}
			testo.EqualValues(t, ok, tt.ok)
		})
	}
}

//func TestAnyValueMatchesTerm(t *testing.T) {
//	tests := []struct {
//		term     string
//		values   []string
//		anyCase  bool
//		contains bool
//		ok       bool
//	}{
//		{"test", []string{"nomatch", "prefixTEST"}, false, false, false},
//		{"prefixTEST", []string{"nomatch", "prefixTEST"}, false, false, true},
//		{"test", []string{"nomatch", "prefixTEST"}, true, false, false},
//		{"prefixtest", []string{"nomatch", "prefixTEST"}, true, false, true},
//		{"test", []string{"nomatch", "prefixTEST"}, false, true, false},
//		{"test", []string{"nomatch", "prefixTEST"}, true, true, true},
//		{"test", []string{"nomatch"}, true, true, false},
//	}
//
//	for ii, tt := range tests {
//		t.Run(strconv.Itoa(ii), func(t *testing.T) {
//			ok := anyValueMatchesTerm(tt.term, tt.values, tt.anyCase, tt.contains)
//			testo.EqualValues(t, ok, tt.ok)
//		})
//	}
//}
//
//func TestReduxMatch(t *testing.T) {
//	tests := []struct {
//		terms    []string
//		scope    map[string]bool
//		anyCase  bool
//		contains bool
//		matches  []string
//	}{
//		{[]string{"11"}, nil, false, false, []string{}},
//		{[]string{"11"}, nil, false, true, []string{"k2"}},
//		{[]string{"11"}, nil, true, true, []string{"k2"}},
//		{[]string{"11"}, map[string]bool{"k1": true, "k3": true}, true, true, []string{}},
//		{[]string{"V12"}, nil, false, false, []string{}},
//		{[]string{"V12"}, nil, true, false, []string{"k2"}},
//		{[]string{"V12"}, nil, false, true, []string{}},
//		{[]string{"V12"}, nil, true, true, []string{"k2", "k3", "k4"}},
//		{[]string{"V12"}, map[string]bool{"k4": true, "k5": true}, true, true, []string{"k4"}},
//	}
//
//	rdx := mockRedux()
//
//	for ii, tt := range tests {
//		t.Run(strconv.Itoa(ii), func(t *testing.T) {
//
//			matches := rdx.Match(tt.terms, tt.scope, tt.anyCase, tt.contains)
//			testo.EqualValues(t, len(matches), len(tt.matches))
//			for _, m := range tt.matches {
//				_, ok := matches[m]
//				testo.EqualValues(t, ok, true)
//			}
//		})
//	}
//}
