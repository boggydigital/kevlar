package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func mockIndex() index {
	idx := index{
		"1": {"1", 1, 1},
		"2": {"2", 2, 2},
		"3": {"3", 3, 3},
	}

	return idx
}

func cleanup() error {
	idxPath := filepath.Join(os.TempDir(), indexFilename)
	if _, err := os.Stat(idxPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return os.Remove(idxPath)
}

func TestIndexReadNothing(t *testing.T) {
	testo.Error(t, cleanup(), false)
	idx := index{}
	testo.Error(t, idx.read(os.TempDir()), false)
	testo.EqualValues(t, len(idx), 0)
}

func TestIndexWriteRead(t *testing.T) {
	idx := mockIndex()
	testo.Error(t, idx.write(os.TempDir()), false)
	idx = index{}
	testo.EqualValues(t, len(idx), 0)
	testo.Error(t, idx.read(os.TempDir()), false)
	testo.DeepEqual(t, idx, mockIndex())
	testo.Error(t, cleanup(), false)
}

func TestIndexUpd(t *testing.T) {
	tests := []struct {
		key         string
		hash        string
		expHash     string
		expModified int
	}{
		{"1", "1", "1", testo.Equal},
		{"1", "2", "2", testo.Greater},
		{"new", "new", "new", testo.Greater},
	}
	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			idx := mockIndex()
			mod := idx[tt.key].Modified

			idx.upd(tt.key, tt.hash)

			testo.EqualValues(t, idx[tt.key].Hash, tt.expHash)
			testo.CompareInt64(t, idx[tt.key].Modified, mod, tt.expModified)
		})
	}
}
