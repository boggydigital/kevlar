package kvas

// reduce filters index records using a provided filter func
func (vs *ValueSet) reduce(filter func(IndexRecord) bool) []string {
	if vs == nil {
		return nil
	}
	keys := make([]string, 0, len(vs.index))
	for k, ir := range vs.index {
		if (filter != nil && filter(ir)) || filter == nil {
			keys = append(keys, k)
		}
	}
	return keys
}

// All returns all value keys of a valueSet
func (vs *ValueSet) All() []string {
	return vs.reduce(nil)
}

// CreatedAfter returns keys of values created on or after provided timestamp
func (vs *ValueSet) CreatedAfter(timestamp int64) []string {
	return vs.reduce(func(ir IndexRecord) bool {
		return ir.Created >= timestamp
	})
}

// ModifiedAfter returns keys of values modified on or after provided timestamp
func (vs *ValueSet) ModifiedAfter(timestamp int64) []string {
	return vs.reduce(func(ir IndexRecord) bool {
		return ir.Modified >= timestamp
	})
}
