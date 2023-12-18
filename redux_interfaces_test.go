package kvas

var (
	_ ReadableRedux  = (*Redux)(nil)
	_ WriteableRedux = (*Redux)(nil)
	_ FixableRedux   = (*Redux)(nil)
)
