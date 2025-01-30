package kevlar_legacy

type LogRecord struct {
	Ts int64
	Mt MutationType
	Id string
}

type LogRecords []*LogRecord
