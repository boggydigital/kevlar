package kvas_compat

const IndexFilename = "_index.gob"

type Record struct {
	Hash     string `json:"hash"`
	Created  int64  `json:"created"`
	Modified int64  `json:"modified"`
}

type Index map[string]Record
