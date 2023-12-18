package kvas

import (
	"github.com/boggydigital/wits"
	"io"
)

func (rdx *Redux) Export(w io.Writer, ids ...string) error {

	skv := make(wits.SectionKeyValues)

	for _, id := range ids {
		skv[id] = make(wits.KeyValues)
		for p := range rdx.assetKeyValues {
			if vals, ok := rdx.GetAllValues(p, id); ok {
				skv[id][p] = vals
			}
		}
	}

	return skv.Write(w)
}
