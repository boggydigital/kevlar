package kevlar_legacy

type logRecord struct {
	Ts int64
	Mt mutationType
	Id string
}

type logRecords []*logRecord
