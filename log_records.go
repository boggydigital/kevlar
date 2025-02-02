package kevlar

const (
	logRecordsFilename = "_log.gob"
)

type logRecord struct {
	Id   string
	Ts   int64
	Mt   MutationType
	Hash []byte
}

type logRecords []*logRecord
