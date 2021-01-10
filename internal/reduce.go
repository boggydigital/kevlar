package internal

// reduce filters index records using a provided filter func
func (vs *valueSet) reduce(filter func(IndexRecord) bool) []string {
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
func (vs *valueSet) All() []string {
	return vs.reduce(nil)
}

// CreatedAfter returns keys of values created on or after provided timestamp
func (vs *valueSet) CreatedAfter(timestamp int64) []string {
	return vs.reduce(func(ir IndexRecord) bool {
		return ir.Created >= timestamp
	})
}

// ModifiedAfter returns keys of values modified on or after provided timestamp
func (vs *valueSet) ModifiedAfter(timestamp int64) []string {
	return vs.reduce(func(ir IndexRecord) bool {
		return ir.Modified >= timestamp
	})
}
