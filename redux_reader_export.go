package kvas

import (
	"github.com/boggydigital/wits"
	"golang.org/x/exp/maps"
	"io"
	"sort"
)

func (rdx *redux) Export(w io.Writer, assets ...string) error {

	if len(assets) == 0 {
		assets = maps.Keys(rdx.assetKeyValues)
	}

	sort.Strings(assets)

	skv := make(wits.SectionKeyValues)

	for _, asset := range assets {
		skv[asset] = rdx.assetKeyValues[asset]
	}

	return skv.Write(w)
}
