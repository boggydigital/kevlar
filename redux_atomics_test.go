package kvas

import (
	"github.com/boggydigital/testo"
	"testing"
)

var mockAtomics = ReduxAtomics{
	"k1": true,
	"k2": true,
}

func TestReduxAtomicsIsAtomic(t *testing.T) {
	tests := []struct {
		key string
		exp bool
	}{
		{"", false},
		{"k1", true},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.EqualValues(t, mockAtomics.IsAtomic(tt.key), tt.exp)
		})
	}
}
