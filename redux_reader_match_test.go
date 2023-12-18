package kvas

import "testing"

func TestRedux_MatchAsset(t *testing.T) {

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
