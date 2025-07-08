package kevlar

import (
	"bytes"
	"github.com/boggydigital/testo"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const (
	testDir = "kevlar_test"
)

func mockKeyValues(t *testing.T) *keyValues {
	return &keyValues{
		dir: filepath.Join(t.TempDir(), testDir),
		ext: GobExt,
		log: []*logRecord{
			{
				Ts: 1,
				Mt: Create,
				Id: "1",
			},
			{
				Ts: 2,
				Mt: Create,
				Id: "2",
			},
			{
				Ts: 3,
				Mt: Update,
				Id: "2",
			},
			{
				Ts: 4,
				Mt: Create,
				Id: "3",
			},
			{
				Ts: 5,
				Mt: Cut,
				Id: "1",
			},
		},
		mtx: new(sync.Mutex),
	}
}

func TestNew(t *testing.T) {
	lkv, err := New(t.TempDir(), JsonExt)
	testo.Nil(t, lkv, false)
	testo.Error(t, err, false)
}

func TestKeyValues_SetHasGetCut(t *testing.T) {
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
			kv, err := New(filepath.Join(t.TempDir(), testDir), GobExt)
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
				err = kv.Cut(ck)
				testo.Error(t, err, false)
			}

		})
	}
}

func TestKeyValues_Since(t *testing.T) {
	tests := []struct {
		since int64
		mts   []MutationType
		exp   map[string]MutationType
	}{
		{-1, []MutationType{}, map[string]MutationType{}},
		{-1, []MutationType{Create}, map[string]MutationType{"1": Create, "2": Create, "3": Create}},
		{-1, []MutationType{Create, Update}, map[string]MutationType{"1": Create, "2": Update, "3": Create}},
		{-1, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut, "2": Update, "3": Create}},

		{0, []MutationType{}, map[string]MutationType{}},
		{0, []MutationType{Create}, map[string]MutationType{"1": Create, "2": Create, "3": Create}},
		{0, []MutationType{Create, Update}, map[string]MutationType{"1": Create, "2": Update, "3": Create}},
		{0, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut, "2": Update, "3": Create}},

		{1, []MutationType{}, map[string]MutationType{}},
		{1, []MutationType{Create}, map[string]MutationType{"1": Create, "2": Create, "3": Create}},
		{1, []MutationType{Create, Update}, map[string]MutationType{"1": Create, "2": Update, "3": Create}},
		{1, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut, "2": Update, "3": Create}},

		{2, []MutationType{}, map[string]MutationType{}},
		{2, []MutationType{Create}, map[string]MutationType{"2": Create, "3": Create}},
		{2, []MutationType{Create, Update}, map[string]MutationType{"2": Update, "3": Create}},
		{2, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut, "2": Update, "3": Create}},

		{3, []MutationType{}, map[string]MutationType{}},
		{3, []MutationType{Create}, map[string]MutationType{"3": Create}},
		{3, []MutationType{Create, Update}, map[string]MutationType{"2": Update, "3": Create}},
		{3, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut, "2": Update, "3": Create}},

		{4, []MutationType{}, map[string]MutationType{}},
		{4, []MutationType{Create}, map[string]MutationType{"3": Create}},
		{4, []MutationType{Create, Update}, map[string]MutationType{"3": Create}},
		{4, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut, "3": Create}},

		{5, []MutationType{}, map[string]MutationType{}},
		{5, []MutationType{Create}, map[string]MutationType{}},
		{5, []MutationType{Create, Update}, map[string]MutationType{}},
		{5, []MutationType{Create, Update, Cut}, map[string]MutationType{"1": Cut}},

		{6, []MutationType{}, map[string]MutationType{}},
		{6, []MutationType{Create}, map[string]MutationType{}},
		{6, []MutationType{Create, Update}, map[string]MutationType{}},
		{6, []MutationType{Create, Update, Cut}, map[string]MutationType{}},
	}

	kv := mockKeyValues(t)

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {

			res := kv.Since(tt.since, tt.mts...)
			testo.Nil(t, res, false)

			for id, mt := range res {
				emt, ok := tt.exp[id]
				testo.EqualValues(t, ok, true)
				testo.EqualValues(t, emt, mt)
			}
		})
	}
}

func TestKeyValues_LogModTime(t *testing.T) {
	start := timeNow()

	kv, err := New(filepath.Join(t.TempDir(), testDir), GobExt)
	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	testo.Error(t, kv.Set("test", strings.NewReader("test")), false)

	cmt := kv.LogModTime("1")
	testo.Error(t, err, false)
	testo.CompareInt64(t, cmt, start, testo.Less)

	cmt = kv.LogModTime("test")
	testo.Error(t, err, false)

	testo.CompareInt64(t, cmt, start, testo.GreaterOrEqual)

	cmt = kv.LogModTime("2")
	testo.Error(t, err, false)
	testo.CompareInt64(t, cmt, start, testo.Less)

	err = kv.Cut("test")
	testo.Error(t, err, false)
}

func TestKeyValues_Len(t *testing.T) {
	kv := mockKeyValues(t)
	testo.Nil(t, kv, false)
	testo.EqualValues(t, kv.Len(), 2) // create 1; create 2; update 2; create 3; cut 1 = [2,3]
}

func TestKeyValues_FileModTime(t *testing.T) {
	start := timeNow()
	kv := mockKeyValues(t)
	testo.Error(t, kv.writeLogRecord(nil), false)

	testo.Error(t, kv.Set("1", strings.NewReader("one")), false)

	mt, err := kv.FileModTime("1")
	testo.Error(t, err, false)
	testo.CompareInt64(t, mt, start, testo.GreaterOrEqual)

	testo.Error(t, kv.Cut("1"), false)

	mt, err = kv.FileModTime("1")
	testo.Error(t, err, false)
	testo.CompareInt64(t, mt, start, testo.Less)
	testo.CompareInt64(t, mt, UnknownModTime, testo.Equal)
}

func TestKeyValues_UpdateCompactsLog(t *testing.T) {
	ikv, err := New(filepath.Join(t.TempDir(), testDir), GobExt)
	testo.Error(t, err, false)

	kv, ok := ikv.(*keyValues)
	testo.EqualValues(t, ok, true)
	testo.Nil(t, kv, false)

	testo.EqualValues(t, len(kv.log), 0)

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // added Create record

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // no writes happened, same content, no new log records

	testo.Error(t, kv.Set("1", strings.NewReader("2")), false)
	testo.EqualValues(t, len(kv.log), 2) // added Update record

	testo.Error(t, kv.Set("1", strings.NewReader("3")), false)
	testo.EqualValues(t, len(kv.log), 2) // existing Update record updated, no new log records

	err = kv.Cut("1")
	testo.Error(t, err, false)
}

func TestKeyValues_CutCompactsLog(t *testing.T) {

	ikv, err := New(filepath.Join(t.TempDir(), testDir), GobExt)
	testo.Error(t, err, false)

	kv, ok := ikv.(*keyValues)
	testo.EqualValues(t, ok, true)
	testo.Nil(t, kv, false)

	testo.EqualValues(t, len(kv.log), 0)

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // added Create record

	testo.Error(t, kv.Set("1", strings.NewReader("1")), false)
	testo.EqualValues(t, len(kv.log), 1) // no writes happened, same content, no new log records

	testo.Error(t, kv.Set("1", strings.NewReader("2")), false)
	testo.EqualValues(t, len(kv.log), 2) // added Update record

	testo.Error(t, kv.Set("1", strings.NewReader("3")), false)
	testo.EqualValues(t, len(kv.log), 2) // existing Update record updated, no new log records

	err = kv.Cut("1")
	testo.EqualValues(t, len(kv.log), 1) // log has been compacted to only store Cut operation
	testo.Error(t, err, false)
}

func TestKeyValues_GoroutineSafe(t *testing.T) {
	kv, err := New(filepath.Join(t.TempDir(), testDir), GobExt)

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
				err := kv.Cut(p + aa)
				testo.Error(t, err, false)
			}
		}(kv, pfx)
	}
	wg.Wait()

	testo.EqualValues(t, kv.Len(), 0)
}

func TestKeyValues_SetNil(t *testing.T) {

	kv, err := New(filepath.Join(t.TempDir(), testDir), GobExt)

	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	err = kv.Set("1", new(bytes.Buffer))
	testo.Error(t, err, false)

	testo.EqualValues(t, kv.Len(), 1)
}
