package kevlar_legacy

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

func TestRedux_Keys(t *testing.T) {
	rdx := mockRedux()
	for asset := range rdx.akv {
		keys := rdx.Keys(asset)
		testo.EqualValues(t, len(keys), len(rdx.akv[asset]))
		for _, k := range keys {
			_, ok := rdx.akv[asset][k]
			testo.EqualValues(t, ok, true)
		}
	}
}

func TestRedux_HasAsset(t *testing.T) {
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
			testo.EqualValues(t, rdx.HasAsset(tt.asset), tt.exp)
		})
	}
}

func TestRedux_HasKey(t *testing.T) {
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

func TestRedux_HasValue(t *testing.T) {
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

func TestRedux_GetAllValues(t *testing.T) {
	tests := []struct {
		asset, key string
		ok         bool
	}{
		{"", "", false},
		{"a1", "k1", true},
		{"1a", "k1", false},
		{"a1", "1k", false},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.asset+tt.key, func(t *testing.T) {
			values, ok := rdx.GetAllValues(tt.asset, tt.key)
			if ok {
				testo.DeepEqual(t, rdx.akv[tt.asset][tt.key], values)
			}
			testo.EqualValues(t, ok, tt.ok)
		})
	}
}

func TestRedux_GetFirstVal(t *testing.T) {
	tests := []struct {
		asset, key string
		ok         bool
	}{
		{"", "", false},
		{"a1", "k1", true},
		{"1a", "k1", false},
		{"a1", "1k", false},
	}

	rdx := mockRedux()
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			fv, ok := rdx.GetLastVal(tt.asset, tt.key)
			if ok {
				testo.DeepEqual(t, rdx.akv[tt.asset][tt.key][0], fv)
			}
			testo.EqualValues(t, ok, tt.ok)
		})
	}
}
