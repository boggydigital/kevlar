package kvas

import (
	"github.com/boggydigital/testo"
	"strconv"
	"testing"
)

var sortableAKV = map[string]map[string][]string{
	"title": {
		"id1": {"A"},
		"id2": {"B"},
		"id3": {"C"},
	},
	"number": {
		"id1": {"2222"},
		"id2": {"3333"},
		"id3": {"1111"},
	},
	"subtitle": {
		"id1": {"Z"},
		"id2": {"Y"},
		"id3": {"X"},
	},
	"binary": {
		"id1": {"true"},
		"id2": {"true"},
		"id3": {"false"},
	},
}

func TestRedux_Sort(t *testing.T) {

	ids := []string{"id1", "id2", "id3"}

	tests := []struct {
		ids    []string
		desc   bool
		assets []string
		exp    []string
		expErr bool
	}{
		{nil, true, nil, []string{}, false},
		{nil, false, nil, []string{}, false},
		{ids, false, []string{"title"}, []string{"id1", "id2", "id3"}, false},
		{ids, true, []string{"title"}, []string{"id3", "id2", "id1"}, false},
		{ids, false, []string{"number"}, []string{"id3", "id1", "id2"}, false},
		{ids, true, []string{"number"}, []string{"id2", "id1", "id3"}, false},
		{ids, false, []string{"binary", "subtitle"}, []string{"id3", "id2", "id1"}, false},
		{ids, true, []string{"binary", "subtitle"}, []string{"id1", "id2", "id3"}, false},
		{nil, false, []string{"asset-that-doesnt-exist"}, nil, true},
		{ids, false, []string{"asset-that-doesnt-exist"}, nil, true},
		{ids, true, []string{"asset-that-doesnt-exist"}, nil, true},
	}

	rdx := &redux{assetKeyValues: sortableAKV}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tr, err := rdx.Sort(tt.ids, tt.desc, tt.assets...)
			testo.Error(t, err, tt.expErr)
			testo.DeepEqual(t, tr, tt.exp)
		})
	}
}
