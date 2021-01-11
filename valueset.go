package kvas

import (
	"github.com/boggydigital/kvas/internal"
)

type ValueSet struct {
	*internal.ValueSet
}

func NewClient(location string, ext string) (*ValueSet, error) {
	jvs, err := internal.NewClient(location, ext)
	return &ValueSet{ValueSet: jvs}, err
}
