package kevlar

func ReduxProxy(idpv map[string]map[string][]string) ReadableRedux {
	assetKeyValues := make(map[string]map[string][]string)
	for id, pv := range idpv {
		for p, v := range pv {
			if assetKeyValues[p] == nil {
				assetKeyValues[p] = make(map[string][]string)
			}
			assetKeyValues[p][id] = v
		}
	}

	return &redux{akv: assetKeyValues}
}
