package kvas

import (
	"github.com/boggydigital/testo"
	"strconv"
	"testing"
)

var mockAggregates = ReduxAggregates{
	"a1": {"v1", "v2", "v3"},
	"a2": {"v4", "v5"},
}

func TestReduxAggregatesIsAggregate(t *testing.T) {
	tests := []struct {
		key string
		exp bool
	}{
		{"", false},
		{"a1", true},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.EqualValues(t, mockAggregates.IsAggregate(tt.key), tt.exp)
		})
	}
}

func TestReduxAggregatesAggregates(t *testing.T) {
	tests := []struct {
		exp map[string]bool
	}{
		{map[string]bool{"a1": true, "a2": true}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			as := mockAggregates.Aggregates()
			testo.EqualValues(t, len(as), len(tt.exp))
			for _, a := range as {
				_, ok := tt.exp[a]
				testo.EqualValues(t, ok, true)
			}
		})
	}
}

func TestReduxAggregatesDetail(t *testing.T) {
	tests := []struct {
		key string
		exp []string
	}{
		{"", nil},
		{"a1", []string{"v1", "v2", "v3"}},
		{"unknown", nil},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.DeepEqual(t, mockAggregates.Detail(tt.key), tt.exp)
		})
	}
}

func TestReduxAggregatesDetailAll(t *testing.T) {
	tests := []struct {
		keys []string
		exp  []string
	}{
		{nil, nil},
		{[]string{""}, []string{""}},
		{[]string{"a1", "a2"}, []string{"v1", "v2", "v3", "v4", "v5"}},
		{[]string{"", "a1"}, []string{"", "v1", "v2", "v3"}},
	}

	for ii, tt := range tests {
		t.Run(strconv.Itoa(ii), func(t *testing.T) {
			details := mockAggregates.DetailAll(tt.keys...)
			testo.EqualValues(t, len(details), len(tt.exp))
			for _, v := range tt.exp {
				_, ok := details[v]
				testo.EqualValues(t, ok, true)
			}
		})
	}
}

func TestReduxAggregatesAggregate(t *testing.T) {
	tests := []struct {
		key string
		exp string
	}{
		{"", ""},
		{"v1", "a1"},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			testo.EqualValues(t, mockAggregates.Aggregate(tt.key), tt.exp)
		})
	}
}
