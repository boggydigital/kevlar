package kvas

import (
	"bytes"
	"github.com/boggydigital/testo"
	"golang.org/x/exp/slices"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func mockLocalKeyValues() *localKeyValues {
	return &localKeyValues{
		idx: mockIndex(),
		mtx: &sync.Mutex{},
	}
}

func TestConnectLocal(t *testing.T) {
	tests := []struct {
		ext    string
		expNil bool
		expErr bool
	}{
		{"", true, true},
		{".txt", true, true},
		{"json", true, true},
		{"gob", true, true},
		{JsonExt, false, false},
		{GobExt, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			lkv, err := ConnectLocal(os.TempDir(), tt.ext)
			testo.Nil(t, lkv, tt.expNil)
			testo.Error(t, err, tt.expErr)

			testo.Error(t, indexCleanup(), false)
		})
	}
}

func TestLocalKeyValuesSetHasGetCut(t *testing.T) {
	tests := []struct {
		set []string
		get map[string]bool
	}{
		{nil, nil},
		{[]string{"x1", "x1"}, map[string]bool{"x1": false}},
		{[]string{"y1", "y2"}, map[string]bool{"y1": false, "y2": false, "y3": true}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			lkv, err := ConnectLocal(os.TempDir(), GobExt)
			testo.Nil(t, lkv, false)
			testo.Error(t, err, false)

			// Set, Has tests
			for _, sk := range tt.set {
				err = lkv.Set(sk, strings.NewReader(sk))
				testo.Error(t, err, false)
				testo.EqualValues(t, lkv.Has(sk), true)
			}

			// Get tests
			for gk, expNil := range tt.get {
				rc, err := lkv.Get(gk)
				testo.Error(t, err, false)
				testo.Nil(t, rc, expNil)

				if expNil {
					continue
				}

				var val []byte
				buf := bytes.NewBuffer(val)
				num, err := io.Copy(buf, rc)
				testo.EqualValues(t, num, int64(len(gk)))
				testo.EqualValues(t, gk, buf.String())

				testo.Error(t, rc.Close(), false)
			}

			// Cut, Has tests

			for _, ck := range tt.set {
				has := lkv.Has(ck)
				ok, err := lkv.Cut(ck)
				testo.EqualValues(t, ok, has)
				testo.Error(t, err, false)
			}

			testo.Error(t, indexCleanup(), false)
		})
	}
}

func TestLocalKeyValues_CreatedAfter(t *testing.T) {

	tests := []struct {
		after int64
		exp   []string
	}{
		{-1, []string{"1", "2", "3"}},
		{0, []string{"1", "2", "3"}},
		{1, []string{"1", "2", "3"}},
		{2, []string{"2", "3"}},
		{3, []string{"3"}},
		{4, []string{}},
	}

	kv := mockLocalKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			ca := kv.CreatedAfter(tt.after)
			testo.EqualValues(t, len(ca), len(tt.exp))
			for _, cav := range ca {
				testo.EqualValues(t, slices.Contains(tt.exp, cav), true)
			}
		})
	}
}

func TestLocalKeyValues_ModifiedAfter(t *testing.T) {

	tests := []struct {
		after int64
		sm    bool
		exp   []string
	}{
		{-1, false, []string{"1", "2", "3"}},
		{0, false, []string{"1", "2", "3"}},
		{1, false, []string{"1", "2", "3"}},
		{2, false, []string{"2", "3"}},
		{3, false, []string{"3"}},
		{4, false, []string{}},
		{-1, true, []string{}},
		{0, true, []string{}},
		{1, true, []string{}},
		{2, true, []string{}},
		{3, true, []string{}},
		{4, true, []string{}},
	}

	kv := mockLocalKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			ma := kv.ModifiedAfter(tt.after, tt.sm)
			testo.EqualValues(t, len(ma), len(tt.exp))
			for _, mav := range ma {
				testo.EqualValues(t, slices.Contains(tt.exp, mav), true)
			}
		})
	}
}

func TestLocalKeyValues_IsModifiedAfter(t *testing.T) {

	tests := []struct {
		key   string
		after int64
		exp   bool
	}{
		{"1", -1, true},
		{"1", 0, true},
		{"1", 1, false},
		{"1", 2, false},
		{"2", 0, true},
		{"2", 1, true},
		{"2", 2, false},
	}

	kv := mockLocalKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			testo.EqualValues(t, kv.IsModifiedAfter(tt.key, tt.after), tt.exp)
		})
	}
}
