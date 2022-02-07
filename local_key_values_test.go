package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"strconv"
	"strings"
	"testing"
)

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

			testo.Error(t, cleanup(), false)
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

				//var val []byte
				//buf := bytes.NewBuffer(val)
				//num, err := io.Copy(buf, rc)
				//testo.EqualValues(t, num, int64(len(gk)))
				//testo.EqualValues(t, gk, buf.String())

				testo.Error(t, rc.Close(), false)
			}

			// Cut, Has tests

			for _, ck := range tt.set {
				has := lkv.Has(ck)
				ok, err := lkv.Cut(ck)
				testo.EqualValues(t, ok, has)
				testo.Error(t, err, false)
			}

			testo.Error(t, cleanup(), false)
		})
	}
}
