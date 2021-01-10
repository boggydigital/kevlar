package kvas

import (
	"github.com/boggydigital/kvas/internal"
)

type JsonValueSet struct {
	*internal.JsonValueSet
}

func NewJsonClient(location string) (*JsonValueSet, error) {
	jvs, err := internal.NewJsonClient(location)
	return &JsonValueSet{JsonValueSet: jvs}, err
}
