package kvas

import (
	"github.com/boggydigital/testo"
	"testing"
)

var mockTransitives = reduxTransitives{
	"from1": "to1",
	"from2": "to2",
}

func TestReduxTransitivesIsTransitive(t *testing.T) {
	tests := []struct {
		key string
		exp bool
	}{
		{"", false},
		{"from1", true},
		{"to1", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.EqualValues(t, mockTransitives.IsTransitive(tt.key), tt.exp)
		})
	}
}

func TestReduxTransitivesTransition(t *testing.T) {
	tests := []struct {
		key string
		exp string
	}{
		{"", ""},
		{"from1", "to1"},
		{"to1", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.EqualValues(t, mockTransitives.Transition(tt.key), tt.exp)
		})
	}
}

func TestReduxTransitivesFmt(t *testing.T) {
	tests := []struct {
		from string
		to   string
		exp  string
	}{
		{"", "", " ()"},
		{"from1", "to1", "to1 (from1)"},
		{"to1", "", " (to1)"},
	}

	for _, tt := range tests {
		t.Run(tt.from+tt.to, func(t *testing.T) {
			testo.EqualValues(t, mockTransitives.Fmt(tt.from, tt.to), tt.exp)
		})
	}
}
