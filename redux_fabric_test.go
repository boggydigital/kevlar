package kvas

import (
	"github.com/boggydigital/testo"
	"strconv"
	"testing"
)

func TestInitFabric(t *testing.T) {
	tests := []struct {
		rf *ReduxFabric
	}{
		{nil},
		{&ReduxFabric{nil, nil, nil}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {

			tt.rf = initFabric(tt.rf)
			testo.Nil(t, tt.rf, false)
			testo.Nil(t, tt.rf.Aggregates, false)
			testo.Nil(t, tt.rf.Transitives, false)
			testo.Nil(t, tt.rf.Atomics, false)

		})
	}
}
