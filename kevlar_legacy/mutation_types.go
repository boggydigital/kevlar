package kevlar_legacy

type mutationType int

const (
	create mutationType = iota
	update
	cut
)
