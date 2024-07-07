package kevlar

type logRecord struct {
	Ts int64
	Mt mutationType
	Id string
}

type logRecords []*logRecord
