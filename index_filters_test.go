package kvas

import (
	"github.com/boggydigital/testo"
	"strconv"
	"sync"
	"testing"
	"time"
)

var mtx = sync.Mutex{}

func TestIndexKeys(t *testing.T) {
	idx := mockIndex()
	keys := idx.Keys(mtx)
	testo.EqualValues(t, len(keys), len(idx))
	for _, k := range keys {
		_, ok := idx[k]
		testo.EqualValues(t, ok, true)
	}
}

func TestIndexCreatedAfter(t *testing.T) {
	tests := []struct {
		new map[string]bool
	}{
		{map[string]bool{"n1": true}},
		{map[string]bool{"n1": true, "n2": true, "n3": true}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			start := time.Now().Unix() - 1
			idx := mockIndex()

			for n := range tt.new {
				idx.upd(n, n)
			}

			res := idx.CreatedAfter(start, mtx)
			testo.EqualValues(t, len(res), len(tt.new))
			for _, r := range res {
				_, ok := tt.new[r]
				testo.EqualValues(t, ok, true)
			}
		})
	}
}

func TestIndexModifiedAfter(t *testing.T) {
	tests := []struct {
		changesNew map[string]bool
	}{
		{},
		{map[string]bool{"1": false, "2": false}},
		{map[string]bool{"1": false, "n1": true}},
		{map[string]bool{"n1": true, "n2": true}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			start := time.Now().Unix() - 1
			idx := mockIndex()

			modOnly, modCre := make(map[string]bool), make(map[string]bool)
			for key, nw := range tt.changesNew {
				if !nw {
					modOnly[key] = true
				}
				modCre[key] = true

				idx.upd(key, "new"+key)
			}

			idxModCre := idx.ModifiedAfter(start, false, mtx)
			testo.EqualValues(t, len(idxModCre), len(modCre))
			for _, k := range idxModCre {
				_, ok := modCre[k]
				testo.EqualValues(t, ok, true)
			}

			idxModOnly := idx.ModifiedAfter(start, true, mtx)
			for _, k := range idxModOnly {
				_, ok := modOnly[k]
				testo.EqualValues(t, ok, true)
			}
		})
	}
}

func TestIndexIsModifiedAfter(t *testing.T) {
	mod := []string{"1", "n1"}
	tests := []struct {
		key string
		exp bool
	}{
		{
			"1",
			true,
		},
		{
			"2",
			false,
		},
		{
			"n1",
			true,
		},
		{
			"n2",
			false,
		},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			start := time.Now().Unix() - 1
			idx := mockIndex()

			for _, k := range mod {
				idx.upd(k, "new"+k)
			}

			ok := idx.IsModifiedAfter(tt.key, start, mtx)
			testo.EqualValues(t, ok, tt.exp)

		})
	}

}
