package kvas

var (
	_ ReadableRedux  = (*redux)(nil)
	_ WriteableRedux = (*redux)(nil)
	_ IndexVetter    = (*redux)(nil)
)
