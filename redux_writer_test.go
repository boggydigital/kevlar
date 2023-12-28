package kvas

import (
	"github.com/boggydigital/testo"
	"os"
	"testing"
)

func TestReduxWriteConnect(t *testing.T) {
	wrdx := mockRedux()

	for asset := range wrdx.assetKeyValues {

		testo.Error(t, wrdx.write(asset), false)

		rdx, err := NewReduxReader(os.TempDir(), asset)
		testo.Error(t, err, false)
		testo.Nil(t, rdx, false)

		testo.Error(t, reduxCleanup(asset), false)
	}
}

//func TestReduxAddVal(t *testing.T) {
//	tests := []struct {
//		key, val string
//		exist    bool
//	}{
//		{"", "new", false},
//		{"k1", "", false},
//		{"k1", "v1", true},
//		{"k1", "v11", false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.key+tt.val, func(t *testing.T) {
//			rdx := mockRedux()
//			testo.EqualValues(t, rdx.HasValue(tt.key, tt.val), tt.exist)
//			testo.Error(t, rdx.AddValues(tt.key, tt.val), false)
//			testo.EqualValues(t, rdx.HasValue(tt.key, tt.val), true)
//			testo.Error(t, reduxCleanup(testAsset), false)
//		})
//	}
//}
//
//func TestReduxReplaceValues(t *testing.T) {
//	tests := []struct {
//		key    string
//		values []string
//	}{
//		{"", []string{"new"}},
//		{"k1", nil},
//		{"k1", []string{""}},
//		{"k1", []string{"v1", "v2"}},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.key+strings.Join(tt.values, ""), func(t *testing.T) {
//			rdx := mockRedux()
//			testo.Error(t, rdx.ReplaceValues(tt.key, tt.values...), false)
//			testo.DeepEqual(t, rdx.keyReductions[tt.key], tt.values)
//			testo.Error(t, reduxCleanup(testAsset), false)
//		})
//	}
//}
//
//func TestReduxBatchReplaceValues(t *testing.T) {
//	tests := []struct {
//		keyValues map[string][]string
//	}{
//		{map[string][]string{"": {"new"}}},
//		{map[string][]string{"k1": nil}},
//		{map[string][]string{"k1": {""}}},
//		{map[string][]string{"k1": {"v1", "v2"}}},
//	}
//
//	for ii, tt := range tests {
//		t.Run(strconv.Itoa(ii), func(t *testing.T) {
//			rdx := mockRedux()
//			testo.Error(t, rdx.BatchReplaceValues(tt.keyValues), false)
//			for key, values := range tt.keyValues {
//				testo.DeepEqual(t, rdx.keyReductions[key], values)
//			}
//			testo.Error(t, reduxCleanup(testAsset), false)
//		})
//	}
//}
//
//func TestReduxCutVal(t *testing.T) {
//	tests := []struct {
//		key, val string
//		exist    bool
//	}{
//		{"", "new", false},
//		{"k1", "", false},
//		{"k1", "v1", true},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.key+tt.val, func(t *testing.T) {
//			rdx := mockRedux()
//			testo.EqualValues(t, rdx.HasValue(tt.key, tt.val), tt.exist)
//			testo.Error(t, rdx.CutVal(tt.key, tt.val), false)
//			testo.EqualValues(t, rdx.HasValue(tt.key, tt.val), false)
//			testo.Error(t, reduxCleanup(testAsset), false)
//		})
//	}
//}
