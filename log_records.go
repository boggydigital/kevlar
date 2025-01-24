package kevlar

type mutationType int

const (
	create mutationType = iota
	update
	cut
)

type logRecord struct {
	Id   string
	Ts   int64
	Mt   mutationType
	Hash []byte
}

type logRecords []*logRecord
