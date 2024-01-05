package kvas

import (
	"github.com/boggydigital/wits"
	"golang.org/x/exp/maps"
	"io"
	"sort"
)

func (rdx *redux) Export(w io.Writer, keys ...string) error {

	assets := maps.Keys(rdx.assetKeyValues)
	sort.Strings(assets)

	skv := make(wits.SectionKeyValues)

	for _, key := range keys {
		skv[key] = make(wits.KeyValues)
		for _, asset := range assets {
			if values := rdx.assetKeyValues[asset][key]; len(values) > 0 {
				skv[key][asset] = values
			}
		}
	}

	return skv.Write(w)
}
