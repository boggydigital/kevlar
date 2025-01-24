package kevlar_legacy

import (
	"github.com/boggydigital/testo"
	"sort"
	"strconv"
	"testing"
)

var matchableAKV = map[string]map[string][]string{
	"t": {
		"1": []string{"title"},
		"2": []string{"title1"},
		"3": []string{"TITLE2"},
	},
	"v": {
		"1": []string{"VALUE1"},
		"2": []string{"value"},
		"3": []string{"value2"},
	},
}

func TestRedux_MatchAsset(t *testing.T) {

	limitedScope := []string{"2", "3"}

	tests := []struct {
		asset   string
		terms   []string
		scope   []string
		options []MatchOption
		exp     []string // expected results should be a-z sorted
	}{
		{"", nil, nil, nil, []string{}},

		{"t", []string{"title"}, nil, nil, []string{"1", "2", "3"}},
		{"t", []string{"title"}, nil, []MatchOption{CaseSensitive}, []string{"1", "2"}},
		{"t", []string{"title"}, nil, []MatchOption{FullMatch}, []string{"1"}},
		{"t", []string{"title"}, nil, []MatchOption{CaseSensitive, FullMatch}, []string{"1"}},

		{"t", []string{"title"}, limitedScope, nil, []string{"2", "3"}},
		{"t", []string{"title"}, limitedScope, []MatchOption{CaseSensitive}, []string{"2"}},
		{"t", []string{"title"}, limitedScope, []MatchOption{FullMatch}, []string{}},
		{"t", []string{"title"}, limitedScope, []MatchOption{CaseSensitive, FullMatch}, []string{}},
	}

	rdx := &redux{akv: matchableAKV}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			found := rdx.MatchAsset(tt.asset, tt.terms, tt.scope, tt.options...)
			// pre-sorting results to avoid comparing same arrays just in different order
			sort.Strings(found)
			testo.DeepEqual(t, found, tt.exp)
		})
	}
}

func TestRedux_Match(t *testing.T) {
	tests := []struct {
		query   map[string][]string
		options []MatchOption
		exp     []string // expected results should be a-z sorted
	}{
		{nil, nil, nil},

		{map[string][]string{"t": {"title"}}, nil, []string{"1", "2", "3"}},
		{map[string][]string{"t": {"title"}}, []MatchOption{CaseSensitive}, []string{"1", "2"}},
		{map[string][]string{"t": {"title"}}, []MatchOption{FullMatch}, []string{"1"}},
		{map[string][]string{"t": {"title"}}, []MatchOption{CaseSensitive, FullMatch}, []string{"1"}},

		{map[string][]string{"v": {"value"}}, nil, []string{"1", "2", "3"}},
		{map[string][]string{"v": {"value"}}, []MatchOption{CaseSensitive}, []string{"2", "3"}},
		{map[string][]string{"v": {"value"}}, []MatchOption{FullMatch}, []string{"2"}},
		{map[string][]string{"v": {"value"}}, []MatchOption{CaseSensitive, FullMatch}, []string{"2"}},

		{map[string][]string{"t": {""}, "v": {"value"}}, nil, []string{"1", "2", "3"}},
		{map[string][]string{"t": {"title"}, "v": {""}}, nil, []string{"1", "2", "3"}},

		{map[string][]string{"t": {"title-that-doesnt-exist"}, "v": {"value"}}, nil, []string{}},
		{map[string][]string{"t": {"title"}, "v": {"value-that-doesnt-exist"}}, nil, []string{}},

		{map[string][]string{"t": {"title"}, "v": {"value"}}, nil, []string{"1", "2", "3"}},
		{map[string][]string{"t": {"title"}, "v": {"value"}}, []MatchOption{CaseSensitive}, []string{"2"}},
		{map[string][]string{"t": {"title"}, "v": {"value"}}, []MatchOption{FullMatch}, []string{}},
	}

	rdx := &redux{akv: matchableAKV}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			found := rdx.Match(tt.query, tt.options...)
			// pre-sorting results to avoid comparing same arrays just in different order
			sort.Strings(found)
			testo.DeepEqual(t, found, tt.exp)
		})
	}
}
