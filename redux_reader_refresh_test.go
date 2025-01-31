package kevlar

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRedux_FileModTime(t *testing.T) {
	start := timeNow()

	wrdx, err := NewReduxWriter(filepath.Join(os.TempDir(), testDir), "test")
	testo.Error(t, err, false)

	rdx := wrdx.(*redux)
	testo.Nil(t, rdx, false)

	// first test: compare unmodified redux mod time
	// expected result: mod time should be less than start of the test

	rmt, err := rdx.FileModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, rmt, start, testo.Less)

	// second test: add a value and compare redux mod time
	// expected result: mod time should be greater or equal than start of the test

	testo.Error(t, rdx.AddValues("test", "k1", "v1"), false)

	rmt, err = rdx.FileModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, rmt, start, testo.GreaterOrEqual)

	// cleanup

	testo.Error(t, rdx.CutValues("test", "k1", "v1"), false)
	ok, err := rdx.kv.Cut("test")
	testo.Error(t, err, false)
	testo.EqualValues(t, ok, true)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestRedux_RefreshReader(t *testing.T) {
	wrdx, err := NewReduxWriter(filepath.Join(os.TempDir(), testDir), "test")
	testo.Error(t, err, false)

	rdx := wrdx.(*redux)
	testo.Nil(t, rdx, false)

	// first test: set modTime to force Refresh and try RefreshReader
	// expected result: redux is refreshed and modTime is updated

	kv, ok := rdx.kv.(*keyValues)
	testo.Nil(t, kv, false)
	testo.EqualValues(t, ok, true)

	testo.Error(t, kv.Set("test", strings.NewReader("test")), false)
	ok, err = kv.Cut("test")
	testo.EqualValues(t, ok, true)
	testo.Error(t, err, false)

	rrdx, err := rdx.RefreshReader()
	testo.Error(t, err, false)

	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	mt, err := rdx.FileModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, mt, -1, testo.Greater)

	// second time: don't change modTime and try to RefreshReader again
	// expected result: no refresh is necessary and modTime is unchanged

	startModTime := mt

	rrdx, err = rdx.RefreshReader()
	testo.Error(t, err, false)

	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	newMt, err := rdx.FileModTime()
	testo.Error(t, err, false)
	testo.EqualValues(t, newMt, startModTime)

	testo.Error(t, logRecordsCleanup(), false)
}
