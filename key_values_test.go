package kevlar

import (
	"bytes"
	"github.com/boggydigital/testo"
	"golang.org/x/exp/slices"
	"io"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
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
		lmt: time.Unix(0, 0),
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
		keys: nil,
		mtx:  new(sync.Mutex),
	}
}

func logRecordsModCleanup() error {
	logModPath := filepath.Join(os.TempDir(), testsDirname, kevlarDirname, logRecordsModFilename)
	if _, err := os.Stat(logModPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return os.Remove(logModPath)
}

func logRecordsCleanup() error {
	logPath := filepath.Join(os.TempDir(), testsDirname, kevlarDirname, logRecordsFilename)
	if _, err := os.Stat(logPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if err := os.Remove(logPath); err != nil {
		return err
	}
	return logRecordsModCleanup()
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
		//{nil, nil},
		//{[]string{"x1", "x1"}, map[string]bool{"x1": false}},
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
				has, err := kv.Has(sk)
				testo.EqualValues(t, has, true)
				testo.Error(t, err, false)
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
				num, err := io.Copy(buf, rc)
				testo.EqualValues(t, num, int64(len(gk)))
				testo.EqualValues(t, gk, buf.String())

				testo.Error(t, rc.Close(), false)
			}

			// Cut, Has tests
			for _, ck := range tt.set {
				has, err := kv.Has(ck)
				testo.Error(t, err, false)
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
			ca, err := kv.CreatedAfter(tt.after)
			testo.Error(t, err, false)
			testo.EqualValues(t, len(ca), len(tt.exp))
			for _, cav := range ca {
				testo.EqualValues(t, slices.Contains(tt.exp, cav), true)
			}
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
			ma, err := kv.UpdatedAfter(tt.after)
			testo.Error(t, err, false)
			testo.EqualValues(t, len(ma), len(tt.exp))
			for _, mav := range ma {
				testo.EqualValues(t, slices.Contains(tt.exp, mav), true)
			}
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
			ma, err := kv.CreatedOrUpdatedAfter(tt.after)
			testo.Error(t, err, false)
			testo.EqualValues(t, len(ma), len(tt.exp))
			for _, mav := range ma {
				testo.EqualValues(t, slices.Contains(tt.exp, mav), true)
			}
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
			ok, err := kv.IsUpdatedAfter(tt.key, tt.after)
			testo.Error(t, err, false)
			testo.EqualValues(t, ok, tt.exp)
		})
	}
}

func TestLocalKeyValues_ModTime(t *testing.T) {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)

	kv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)
	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	testo.Error(t, kv.Set("test", strings.NewReader("test")), false)

	cmt, err := kv.ModTime("1")
	testo.Error(t, err, false)
	testo.EqualValues(t, cmt.Before(start), true)

	cmt, err = kv.ModTime("test")
	testo.Error(t, err, false)

	testo.EqualValues(t, cmt.After(start), true)

	cmt, err = kv.ModTime("2")
	testo.Error(t, err, false)
	testo.EqualValues(t, cmt.Before(start), true)

	ok, err := kv.Cut("test")
	testo.EqualValues(t, ok, true)
	testo.Error(t, err, false)

	testo.Error(t, logRecordsCleanup(), false)
}

func remove5Values(kv KeyValues, pfx string) {
	for ii := 0; ii < 5; ii++ {

		aa := strconv.FormatInt(int64(ii), 10)
		ok, err := kv.Cut(pfx + aa)
		if err != nil {
			log.Println(err)
		}
		if !ok {
			log.Println(pfx+aa, "not found")
		}
		d := time.Duration(rand.N(5)+1) * time.Millisecond
		time.Sleep(d)
	}
}

func TestKeyValues_GoroutineSafe(t *testing.T) {
	kv, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)

	testo.Nil(t, kv, false)
	testo.Error(t, err, false)

	pfxs := []string{"a", "b"}
	vals := 5

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
				if err != nil {
					log.Println(err)
				}
			}
		}(kv, pfx)
	}
	wg.Wait()

	keys, err := kv.Keys()
	testo.Error(t, err, false)
	testo.EqualValues(t, len(keys), len(pfxs)*vals)

	for _, pfx := range pfxs {
		wg.Add(1)
		go func(kv KeyValues, p string) {
			defer wg.Done()
			for ii := 0; ii < vals; ii++ {
				aa := strconv.FormatInt(int64(ii), 10)
				// don't check if the value was removed - this will be validated below
				_, err := kv.Cut(p + aa)
				if err != nil {
					log.Println(err)
				}
			}
		}(kv, pfx)
	}
	wg.Wait()

	keys, err = kv.Keys()
	testo.Error(t, err, false)
	testo.EqualValues(t, len(keys), 0)

	testo.Error(t, logRecordsCleanup(), false)
}

func TestKeyValues_MultiInstanceSafe(t *testing.T) {
	kv1, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)

	testo.Error(t, err, false)
	testo.Nil(t, kv1, false)

	kv2, err := NewKeyValues(filepath.Join(os.TempDir(), testsDirname), GobExt)

	testo.Error(t, err, false)
	testo.Nil(t, kv2, false)

	kvs := []KeyValues{kv1, kv2}
	pfxs := []string{"a", "b"}
	vals := 2

	testo.EqualValues(t, len(kvs), len(pfxs))

	// first: concurrently set values in groups
	// in the end keyValues should contain all values

	for pp, pfx := range pfxs {
		func(kv KeyValues, p string) {
			for ii := 0; ii < vals; ii++ {
				aa := strconv.FormatInt(int64(ii), 10)
				err := kv.Set(p+aa, strings.NewReader(aa))
				if err != nil {
					log.Println(err)
				}
			}
		}(kvs[pp], pfx)
	}

	keys1, err := kv1.Keys()
	testo.Error(t, err, false)
	keys2, err := kv2.Keys()
	testo.Error(t, err, false)

	testo.EqualValues(t, len(keys1), len(keys2))
	testo.EqualValues(t, len(keys1), len(pfxs)*vals)

	for pp, pfx := range pfxs {
		func(kv KeyValues, p string) {
			for ii := 0; ii < vals; ii++ {
				aa := strconv.FormatInt(int64(ii), 10)
				// don't check if the value was removed - this will be validated below
				_, err := kv.Cut(p + aa)
				if err != nil {
					log.Println(err)
				}
			}
		}(kvs[pp], pfx)
	}

	keys1, err = kv1.Keys()
	testo.Error(t, err, false)
	keys2, err = kv2.Keys()
	testo.Error(t, err, false)

	testo.EqualValues(t, len(keys1), len(keys2))
	testo.EqualValues(t, len(keys1), 0)

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
