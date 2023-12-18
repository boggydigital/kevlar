package kvas

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
	for p, _ := range is.properties {
		if is.ipv[i].values[p] == is.ipv[j].values[p] {
			continue
		}
		return is.ipv[i].values[p] < is.ipv[j].values[p]
	}
	return false
}
