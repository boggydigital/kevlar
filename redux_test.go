package kevlar

import (
	"os"
	"path/filepath"
	"sync"
)

func reduxCleanup(assets ...string) error {
	for _, asset := range assets {
		rdxPath := filepath.Join(os.TempDir(), testDir, asset+GobExt)
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
	return logRecordsCleanup()
}

func mockRedux() *redux {
	return &redux{
		dir: filepath.Join(os.TempDir(), testDir),
		akv: map[string]map[string][]string{
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
		kv:  mockKeyValues(),
		mtx: new(sync.Mutex),
	}
}
