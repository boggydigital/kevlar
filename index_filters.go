package kvas

import "sync"

// filter returns a set of keys matching a provided filter func
func (idx index) filter(f func(record) bool, mtx *sync.Mutex) []string {

	keys := make([]string, 0, len(idx))

	mtx.Lock()
	defer mtx.Unlock()

	for k, ir := range idx {
		if (f != nil && f(ir)) || f == nil {
			keys = append(keys, k)
		}
	}

	return keys
}

// Keys returns all keys of a keyValues
func (idx index) Keys(mtx *sync.Mutex) []string {
	return idx.filter(nil, mtx)
}

// CreatedAfter returns keys of values created on or after provided timestamp
func (idx index) CreatedAfter(timestamp int64, mtx *sync.Mutex) []string {
	return idx.filter(func(ir record) bool {
		return ir.Created >= timestamp
	}, mtx)
}

// ModifiedAfter returns keys of values modified on or after provided timestamp
// that were created earlier
func (idx index) ModifiedAfter(timestamp int64, strictlyModified bool, mtx *sync.Mutex) []string {
	return idx.filter(func(ir record) bool {
		if strictlyModified && ir.Modified == ir.Created {
			return false
		}
		return ir.Modified >= timestamp
	}, mtx)
}

func (idx index) IsModifiedAfter(key string, timestamp int64, mtx *sync.Mutex) bool {
	mtx.Lock()
	defer mtx.Unlock()

	if ir, ok := idx[key]; ok {
		return ir.Modified > timestamp
	}
	return false
}
