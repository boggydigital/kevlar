package kevlar

import (
	"encoding/gob"
	"errors"
	"github.com/boggydigital/kevlar/kvas_compat"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Migrate transforms kvas index to kevlar log and hash files:
// 0) not part of this method - make sure to create backups before migration!
// 1) existing index file is decoded to get created, modified dates and hashes
// 2) new keyValues is connected to the same directory and cast to specific type
// to get access to internal methods and types (e.g. logRecords)
// 3) every index record is translated to corresponding log record values
// 4) hash files are created for each index record with hash
// 5) log is written as a single operation (vs kv.appendLogRecord calls)
// 6) old index is removed to make sure calling migrate again doesn't overwrite new data
func Migrate(dir string) error {

	// 1)

	absIndexFilename := filepath.Join(dir, kvas_compat.IndexFilename)

	if _, err := os.Stat(absIndexFilename); os.IsNotExist(err) {
		// if index file doesn't exist - don't throw error
		// assuming the migration already happened and there's
		// nothing else to do
		return nil
	} else if err != nil {
		return err
	}

	indexFile, err := os.Open(absIndexFilename)
	if err != nil {
		return err
	}
	defer indexFile.Close()

	var index kvas_compat.Index

	if err = gob.NewDecoder(indexFile).Decode(&index); err != nil {
		return err
	}

	// 2)

	// we won't be writing anything that requires extension, so it
	// can safely be set to an empty string
	ikv, err := NewKeyValues(dir, "")
	if err != nil {
		return err
	}

	kv, ok := ikv.(*keyValues)
	if !ok {
		return errors.New("kevlar: unable to cast interface to a specific type")
	}

	// 3)

	for id, indexRecord := range index {

		kv.log = append(kv.log, &logRecord{
			Ts: indexRecord.Created,
			Mt: create,
			Id: id,
		})

		if indexRecord.Modified > indexRecord.Created {
			kv.log = append(kv.log, &logRecord{
				Ts: indexRecord.Modified,
				Mt: update,
				Id: id,
			})
		}

		// 4)

		if err = kv.createHashFile(id, indexRecord.Hash); err != nil {
			return err
		}
	}

	// 5)

	logRecordsFile, err := os.Create(kv.absLogRecordsFilename())
	if err != nil {
		return err
	}
	defer logRecordsFile.Close()

	if err = gob.NewEncoder(logRecordsFile).Encode(kv.log); err != nil {
		return err
	}

	// 6)

	if err = os.Remove(absIndexFilename); err != nil {
		return err
	}

	return nil
}

// MigrateAll looks for index files in the provided directory and each
// subdirectory and migrates every key values store that is found
func MigrateAll(dir string) error {

	matches := make([]string, 0)
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, kvas_compat.IndexFilename) {
			matches = append(matches, path)
			return nil
		}
		return nil
	}); err != nil {
		return err
	}

	for _, match := range matches {
		indexDir, _ := filepath.Split(match)
		if err := Migrate(indexDir); err != nil {
			return err
		}
	}

	return nil
}
