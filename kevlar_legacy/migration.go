package kevlar_legacy

const (
	KevlarDirname = "_kevlar"
	HashExt       = ".sha256"
)

type MutationType int

const (
	Create MutationType = iota
	Update
	Cut
)

type LogRecord struct {
	Ts int64
	Mt MutationType
	Id string
}

type LogRecords []*LogRecord
