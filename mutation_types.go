package kevlar

type mutationType int

const (
	create mutationType = iota
	update
	cut
)
