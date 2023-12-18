package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"path/filepath"
	"testing"
)

const (
	testAsset   = "test_asset"
	detailAsset = "detail_asset"
)

func reduxCleanup(assets ...string) error {
	for _, asset := range assets {
		rdxPath := filepath.Join(os.TempDir(), asset+GobExt)
		if _, err := os.Stat(rdxPath); err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if err := os.Remove(rdxPath); err != nil {
			return err
		}
	}

	return indexCleanup()
}

func mockRedux() *Redux {
	return &Redux{
		dir: os.TempDir(),
		assetKeyValues: map[string]map[string][]string{
			"a1": {
				"k1": {"v11"},
				"k2": {"v21", "v22"},
				"k3": {"v31", "v32", "v33"},
			},
			"a2": {
				"k4": {"v41", "v42", "v43", "v44"},
				"k5": {"v51", "v52", "v53", "v54", "v55"},
			},
		},
	}
}

func TestReduxWriteConnect(t *testing.T) {
	wrdx := mockRedux()

	testo.Error(t, wrdx.write("a1"), false)

	rdx, err := ReadRedux(os.TempDir(), "a1")
	testo.Error(t, err, false)
	testo.Nil(t, rdx, false)

	testo.Error(t, reduxCleanup("a1"), false)
}
