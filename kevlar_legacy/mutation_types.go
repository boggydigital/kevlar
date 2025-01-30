package kevlar_legacy

type MutationType int

const (
	Create MutationType = iota
	Update
	Cut
)
