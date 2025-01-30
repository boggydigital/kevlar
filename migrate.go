package kevlar

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/boggydigital/kevlar/kevlar_legacy"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func Migrate(dir string, ext string) error {

	// load legacy log
	absLegacyLogFilename := filepath.Join(dir, kevlar_legacy.KevlarDirname, logRecordsFilename)
	if _, err := os.Stat(absLegacyLogFilename); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	logFile, err := os.Open(absLegacyLogFilename)
	if err != nil {
		return err
	}
	defer logFile.Close()

	var legacyLog kevlar_legacy.LogRecords

	if err := gob.NewDecoder(logFile).Decode(&legacyLog); err == io.EOF {
		// do nothing - empty log will be initialized later
	} else if err != nil {
		return err
	}

	// load legacy hashes
	legacyHashes := make(map[string][]byte)
	for _, lr := range legacyLog {
		id := lr.Id
		hash, err := readLegacyHash(id, dir)
		if err != nil {
			return err
		}
		legacyHashes[id] = hash
	}

	// recreate new log

	var newLogRecords logRecords

	for _, llr := range legacyLog {

		newLr := &logRecord{
			Id:   llr.Id,
			Ts:   llr.Ts,
			Mt:   mapMutationType(llr.Mt),
			Hash: legacyHashes[llr.Id],
		}

		newLogRecords = append(newLogRecords, newLr)
	}

	kv := &keyValues{
		dir: dir,
		ext: ext,
		log: newLogRecords,
		mtx: new(sync.Mutex),
	}

	if err := kv.writeLogRecord(nil); err != nil {
		return err
	}

	// delete legacy filed (log, hashes)

	for id := range legacyHashes {
		hashFilename := filepath.Join(dir, kevlar_legacy.KevlarDirname, id+kevlar_legacy.HashExt)
		if _, err := os.Stat(hashFilename); os.IsNotExist(err) {
			continue
		}
		if err := os.Remove(hashFilename); err != nil {
			return err
		}
	}

	if err := os.Remove(absLegacyLogFilename); err != nil {
		return err
	}

	return os.Remove(filepath.Join(dir, kevlar_legacy.KevlarDirname))
}

func readLegacyHash(id, dir string) ([]byte, error) {

	hashFilename := filepath.Join(dir, kevlar_legacy.KevlarDirname, id+kevlar_legacy.HashExt)

	if _, err := os.Stat(hashFilename); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	hashFile, err := os.Open(hashFilename)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, hashFile); err != nil {
		return nil, err
	}

	return hex.DecodeString(buf.String())
}

func hashFile(id, dir, ext string) ([]byte, error) {
	filename := filepath.Join(dir, id+ext)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return sha256Bytes(file)
}

func mapMutationType(mt kevlar_legacy.MutationType) mutationType {
	switch mt {
	case kevlar_legacy.Create:
		return create
	case kevlar_legacy.Update:
		return update
	case kevlar_legacy.Cut:
		return cut
	}
	return create
}
