package kevlar

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRedux_ModTime(t *testing.T) {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)

	wrdx, err := NewReduxWriter(filepath.Join(os.TempDir(), testsDirname), "test")
	testo.Error(t, err, false)

	rdx := wrdx.(*redux)
	testo.Nil(t, rdx, false)

	// first test: compare unmodified redux mod time
	// expected result: mod time should be less than start of the test

	rmt, err := rdx.ModTime()
	testo.Error(t, err, false)
	testo.EqualValues(t, start.After(rmt), true)

	// second test: add a value and compare redux mod time
	// expected result: mod time should be greater or equal than start of the test

	testo.Error(t, rdx.AddValues("test", "k1", "v1"), false)

	rmt, err = rdx.ModTime()
	testo.Error(t, err, false)
	testo.EqualValues(t, rmt.After(start), true)

	// cleanup

	testo.Error(t, rdx.CutValues("test", "k1", "v1"), false)
	ok, err := rdx.kv.Cut("test")
	testo.Error(t, err, false)
	testo.EqualValues(t, ok, true)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestRedux_RefreshReader(t *testing.T) {
	wrdx, err := NewReduxWriter(filepath.Join(os.TempDir(), testsDirname), "test")
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

	testo.EqualValues(t, rdx.lmt.After(time.Unix(0, 0)), true)

	// second time: don't change modTime and try to RefreshReader again
	// expected result: no refresh is necessary and modTime is unchanged

	startModTime := rdx.lmt

	rrdx, err = rdx.RefreshReader()
	testo.Error(t, err, false)

	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	testo.EqualValues(t, rdx.lmt, startModTime)

	testo.Error(t, logRecordsCleanup(), false)
}
