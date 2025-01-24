package kevlar

import (
	"bytes"
	"github.com/boggydigital/testo"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	testsDirname = "kevlar_tests"
)

func mockKeyValues() *keyValues {
	return &keyValues{
		dir: filepath.Join(os.TempDir(), testsDirname),
		ext: GobExt,
		log: []*logRecord{
			{
				Ts: 1,
				Mt: create,
				Id: "1",
			},
			{
				Ts: 2,
				Mt: create,
				Id: "2",
			},
			{
				Ts: 3,
				Mt: update,
				Id: "2",
			},
			{
				Ts: 4,
				Mt: create,
				Id: "3",
			},
			{
				Ts: 5,
				Mt: cut,
				Id: "1",
			},
		},
		mtx: new(sync.Mutex),
	}
}

func logRecordsCleanup() error {
	logPath := filepath.Join(os.TempDir(), testsDirname, logRecordsFilename)
	if _, err := os.Stat(logPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.Remove(logPath); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(os.TempDir(), testsDirname))
}

func TestNewKeyValues(t *testing.T) {
	lkv, err := NewKeyValues(os.TempDir(), JsonExt)
	testo.Nil(t, lkv, false)
	testo.Error(t, err, false)
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
			kv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)
			testo.Nil(t, kv, false)
			testo.Error(t, err, false)

			// Set, Has tests
			for _, sk := range tt.set {
				err := kv.Set(sk, strings.NewReader(sk))
				testo.Error(t, err, false)
				has := kv.Has(sk)
				testo.EqualValues(t, has, true)
			}

			// Get tests
			for gk, expNil := range tt.get {
				rc, err := kv.Get(gk)
				testo.Error(t, err, expNil)
				testo.Nil(t, rc, expNil)

				if expNil {
					continue
				}

				var val []byte
				buf := bytes.NewBuffer(val)
				var cerr error
				num, cerr := io.Copy(buf, rc)
				testo.Error(t, cerr, false)
				testo.EqualValues(t, num, int64(len(gk)))
				testo.EqualValues(t, gk, buf.String())

				testo.Error(t, rc.Close(), false)
			}

			// Cut, Has tests
			for _, ck := range tt.set {
				has := kv.Has(ck)
				ok, err := kv.Cut(ck)
				testo.EqualValues(t, ok, has)
				testo.Error(t, err, false)
			}

			testo.Error(t, logRecordsCleanup(), false)

		})
	}
}

func TestLocalKeyValues_CreatedAfter(t *testing.T) {

	tests := []struct {
		after int64
		exp   []string
	}{
		{-1, []string{"2", "3"}},
		{0, []string{"2", "3"}},
		{1, []string{"2", "3"}},
		{2, []string{"2", "3"}},
		{3, []string{"3"}},
		{4, []string{"3"}},
		{5, []string{}},
	}

	kv := mockKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			ca := kv.CreatedAfter(tt.after)
			caLen := 0
			for cav := range ca {
				caLen++
				testo.EqualValues(t, slices.Contains(tt.exp, cav), true)
			}
			testo.EqualValues(t, caLen, len(tt.exp))
		})
	}
}

func TestLocalKeyValues_UpdatedAfter(t *testing.T) {

	tests := []struct {
		after int64
		exp   []string
	}{
		{-1, []string{"2"}},
		{0, []string{"2"}},
		{1, []string{"2"}},
		{2, []string{"2"}},
		{3, []string{"2"}},
		{4, []string{}},
	}

	kv := mockKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			ua := kv.UpdatedAfter(tt.after)
			uaLen := 0
			for uav := range ua {
				uaLen++
				testo.EqualValues(t, slices.Contains(tt.exp, uav), true)
			}
			testo.EqualValues(t, uaLen, len(tt.exp))
		})
	}
}

func TestLocalKeyValues_CreatedOrUpdatedAfter(t *testing.T) {

	tests := []struct {
		after int64
		exp   []string
	}{
		{-1, []string{"2", "3"}},
		{0, []string{"2", "3"}},
		{1, []string{"2", "3"}},
		{2, []string{"2", "3"}},
		{3, []string{"2", "3"}},
		{4, []string{"3"}},
		{5, []string{}},
	}

	kv := mockKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			cua := kv.CreatedOrUpdatedAfter(tt.after)
			cuaLen := 0
			for cuav := range cua {
				cuaLen++
				testo.EqualValues(t, slices.Contains(tt.exp, cuav), true)
			}
			testo.EqualValues(t, cuaLen, len(tt.exp))
		})
	}
}

func TestLocalKeyValues_IsUpdatedAfter(t *testing.T) {

	tests := []struct {
		key   string
		after int64
		exp   bool
	}{
		{"1", -1, false},
		{"1", 0, false},
		{"1", 1, false},
		{"1", 2, false},
		{"2", 0, true},
		{"2", 1, true},
		{"2", 2, true},
		{"2", 3, true},
		{"2", 4, false},
		{"3", -1, false},
	}

	kv := mockKeyValues()
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			ok := kv.IsUpdatedAfter(tt.key, tt.after)
			testo.EqualValues(t, ok, tt.exp)
		})
	}
}

func TestLocalKeyValues_ValueModTime(t *testing.T) {
	start := time.Now().Unix()
	time.Sleep(100 * time.Millisecond)

	kv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)
	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	testo.Error(t, kv.Set("test", strings.NewReader("test")), false)

	cmt := kv.ValueModTime("1")
	testo.Error(t, err, false)
	testo.CompareInt64(t, cmt, start, testo.Less)

	cmt = kv.ValueModTime("test")
	testo.Error(t, err, false)

	testo.CompareInt64(t, cmt, start, testo.GreaterOrEqual)

	cmt = kv.ValueModTime("2")
	testo.Error(t, err, false)
	testo.CompareInt64(t, cmt, start, testo.Less)

	ok, err := kv.Cut("test")
	testo.EqualValues(t, ok, true)
	testo.Error(t, err, false)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestKeyValues_UpdatesPreventLogGrowth(t *testing.T) {
	ikv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)
	testo.Error(t, err, false)

	kv, ok := ikv.(*keyValues)
	testo.EqualValues(t, ok, true)
	testo.Nil(t, kv, false)

	testo.EqualValues(t, len(kv.log), 0)

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // added create record

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // no writes happened, same content, no new log records

	testo.Error(t, kv.Set("1", strings.NewReader("2")), false)
	testo.EqualValues(t, len(kv.log), 2) // added update record

	testo.Error(t, kv.Set("1", strings.NewReader("3")), false)
	testo.EqualValues(t, len(kv.log), 2) // existing update record updated, no new log records

	ok, err = kv.Cut("1")
	testo.EqualValues(t, ok, true)
	testo.Error(t, err, false)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestKeyValues_CutCompactsLog(t *testing.T) {

	ikv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)
	testo.Error(t, err, false)

	kv, ok := ikv.(*keyValues)
	testo.EqualValues(t, ok, true)
	testo.Nil(t, kv, false)

	testo.EqualValues(t, len(kv.log), 0)

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // added create record

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // no writes happened, same content, no new log records

	testo.Error(t, kv.Set("1", strings.NewReader("2")), false)
	testo.EqualValues(t, len(kv.log), 2) // added update record

	testo.Error(t, kv.Set("1", strings.NewReader("3")), false)
	testo.EqualValues(t, len(kv.log), 2) // existing update record updated, no new log records

	ok, err = kv.Cut("1")
	testo.EqualValues(t, len(kv.log), 1) // log has been compacted to only store cut operation
	testo.EqualValues(t, ok, true)
	testo.Error(t, err, false)

	testo.Error(t, logRecordsCleanup(), false)

}

func TestKeyValues_GoroutineSafe(t *testing.T) {
	kv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)

	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	pfxs := []string{"a", "b", "c", "d", "e"}
	vals := 10

	// first: concurrently set values in groups
	// in the end keyValues should contain all values

	var wg sync.WaitGroup
	for _, pfx := range pfxs {
		wg.Add(1)
		go func(kv KeyValues, p string) {
			defer wg.Done()
			for ii := 0; ii < vals; ii++ {
				aa := strconv.FormatInt(int64(ii), 10)
				err := kv.Set(p+aa, strings.NewReader(aa))
				testo.Error(t, err, false)
			}
		}(kv, pfx)
	}
	wg.Wait()

	testo.EqualValues(t, kv.Len(), len(pfxs)*vals)

	for _, pfx := range pfxs {
		wg.Add(1)
		go func(kv KeyValues, p string) {
			defer wg.Done()
			for ii := 0; ii < vals; ii++ {
				aa := strconv.FormatInt(int64(ii), 10)
				ok, err := kv.Cut(p + aa)
				testo.EqualValues(t, ok, true)
				testo.Error(t, err, false)
			}
		}(kv, pfx)
	}
	wg.Wait()

	testo.EqualValues(t, kv.Len(), 0)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestKeyValues_NoExternalModificationsAllowed(t *testing.T) {

	kv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)
	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	err = kv.Set("1", strings.NewReader("1"))
	testo.Error(t, err, false)

	time.Sleep(time.Millisecond * 1000)

	kvLog, err := os.Create(filepath.Join(os.TempDir(), testsDirname, logRecordsFilename))
	testo.Error(t, err, false)

	_, err = io.WriteString(kvLog, "")
	testo.Error(t, err, false)

	// this should return error as the log has been modified externally
	err = kv.Set("1", strings.NewReader("2"))
	testo.Error(t, err, true)

	testo.Error(t, logRecordsCleanup(), false)
}
