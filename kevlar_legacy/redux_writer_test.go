package kevlar_legacy

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReduxWriteConnect(t *testing.T) {
	wrdx := mockRedux()

	for asset := range wrdx.akv {

		testo.Error(t, wrdx.write(asset), false)

		rdx, err := NewReduxReader(filepath.Join(os.TempDir(), testsDirname), asset)
		testo.Error(t, err, false)
		testo.Nil(t, rdx, false)
		testo.Error(t, rdx.MustHave(asset), false)

		testo.Error(t, reduxCleanup(asset), false)
	}
}

func TestReduxAddVal(t *testing.T) {
	tests := []struct {
		asset, key, val string
		exist           bool
	}{
		{"a1", "", "new", false},
		{"a1", "k1", "", false},
		{"a1", "k1", "v11", true},
		{"a1", "k1", "v1", false},
	}

	for _, tt := range tests {
		t.Run(tt.key+tt.val, func(t *testing.T) {
			rdx := mockRedux()
			testo.EqualValues(t, rdx.HasValue(tt.asset, tt.key, tt.val), tt.exist)
			testo.Error(t, rdx.AddValues(tt.asset, tt.key, tt.val), false)
			testo.EqualValues(t, rdx.HasValue(tt.asset, tt.key, tt.val), true)
			// cleanup
			testo.Error(t, rdx.CutValues(tt.asset, tt.key, tt.val), false)
			ok, err := rdx.kv.Cut("a1")
			testo.Error(t, err, false)
			testo.EqualValues(t, ok, true)
			testo.Error(t, logRecordsCleanup(), false)
		})
	}
}

func TestReduxReplaceValues(t *testing.T) {
	tests := []struct {
		asset, key string
		values     []string
	}{
		{"a1", "k1", nil},
		{"a1", "k1", []string{""}},
		{"a1", "k1", []string{"v1", "v2"}},
	}

	for _, tt := range tests {
		t.Run(tt.key+strings.Join(tt.values, ""), func(t *testing.T) {
			rdx := mockRedux()
			testo.Error(t, rdx.ReplaceValues(tt.asset, tt.key, tt.values...), false)
			testo.DeepEqual(t, rdx.akv[tt.asset][tt.key], tt.values)
			// cleanup
			testo.Error(t, rdx.CutValues(tt.asset, tt.key, tt.values...), false)
			ok, err := rdx.kv.Cut("a1")
			testo.Error(t, err, false)
			testo.EqualValues(t, ok, true)
			testo.Error(t, logRecordsCleanup(), false)
		})
	}
}
