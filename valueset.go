package kvas

import (
	"fmt"
	"github.com/boggydigital/kvas/internal"
	"strings"
)

type ValueSet struct {
	*internal.ValueSet
}

func NewClient(location string, ext string, debug bool) (*ValueSet, error) {
	if !strings.HasPrefix(ext, ".") {
		return nil, fmt.Errorf("extension should start with a '.', e.g. '.ext'")
	}
	jvs, err := internal.NewClient(location, ext, debug)
	return &ValueSet{ValueSet: jvs}, err
}
