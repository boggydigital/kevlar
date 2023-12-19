package kvas

import (
	"github.com/boggydigital/testo"
	"testing"
	"time"
)

func TestRedux_ModTime(t *testing.T) {
	start := time.Now().Unix()

	rdx := mockRedux()

	// first test: compare unmodified redux mod time
	// expected result: mod time should be less than start of the test

	rmt, err := rdx.ModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, rmt, start, testo.Less)

	// second test: add a value and compare redux mod time
	// expected result: mod time should be greater or equal than start of the test

	testo.Error(t, rdx.AddValues("a1", "k1", "test"), false)

	rmt, err = rdx.ModTime()
	testo.Error(t, err, false)
	testo.CompareInt64(t, rmt, start, testo.GreaterOrEqual)

	testo.Error(t, cleanupLocalKeyValues(rdx.kv.(*localKeyValues)), false)
}

func TestRedux_RefreshReader(t *testing.T) {
	rdx := mockRedux()

	// first test: set modTime to force Refresh and try RefreshReader
	// expected result: redux is refreshed and modTime is updated

	rdx.modTime = -1
	rrdx, err := rdx.RefreshReader()
	testo.Error(t, err, false)

	var ok bool
	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	testo.CompareInt64(t, rdx.modTime, -1, testo.Greater)

	// second time: don't change modTime and try to RefreshReader again
	// expected result: no refresh is necessary and modTime is unchanged

	startModTime := rdx.modTime

	rrdx, err = rdx.RefreshReader()
	testo.Error(t, err, false)

	rdx, ok = rrdx.(*redux)
	testo.EqualValues(t, ok, true)

	testo.EqualValues(t, rdx.modTime, startModTime)
}
