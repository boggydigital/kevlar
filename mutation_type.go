package kevlar

type MutationType int

const (
	Create MutationType = iota
	Update
	Cut
)

var mutationTypeStrings = map[MutationType]string{
	Create: "create",
	Update: "update",
	Cut:    "cut",
}

func (mt MutationType) String() string {
	return mutationTypeStrings[mt]
}
