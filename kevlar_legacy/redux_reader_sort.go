package kevlar_legacy

import "sort"

type idValues struct {
	id     string
	values []string
}

type sortableIdSet struct {
	properties []string
	ipv        []idValues
}

func (is *sortableIdSet) Len() int {
	return len(is.ipv)
}

func (is *sortableIdSet) Swap(i, j int) {
	is.ipv[i], is.ipv[j] = is.ipv[j], is.ipv[i]
}

func (is *sortableIdSet) Less(i, j int) bool {
	for p := range is.properties {
		if is.ipv[i].values[p] == is.ipv[j].values[p] {
			continue
		}
		return is.ipv[i].values[p] < is.ipv[j].values[p]
	}
	return false
}

func (rdx *redux) Sort(ids []string, desc bool, sortBy ...string) ([]string, error) {
	if err := rdx.MustHave(sortBy...); err != nil {
		return nil, err
	}

	sis := &sortableIdSet{
		properties: sortBy,
		ipv:        make([]idValues, 0, len(ids)),
	}

	for _, id := range ids {
		iv := idValues{id: id}
		for _, p := range sortBy {
			v, _ := rdx.GetLastVal(p, id)
			iv.values = append(iv.values, v)
		}
		sis.ipv = append(sis.ipv, iv)
	}

	var sortInterface sort.Interface = sis
	if desc {
		sortInterface = sort.Reverse(sortInterface)
	}

	sort.Sort(sortInterface)

	sorted := make([]string, 0, len(sis.ipv))
	for _, iv := range sis.ipv {
		sorted = append(sorted, iv.id)
	}

	return sorted, nil
}
