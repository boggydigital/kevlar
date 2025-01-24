package main

import (
	"fmt"
	"github.com/boggydigital/kevlar"
	"github.com/boggydigital/kevlar/kevlar_legacy"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	start := time.Now().Unix()
	fmt.Println("start:", start)

	uhd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	downloadsDir := filepath.Join(uhd, "Downloads")
	key := "1"

	kv, err := kevlar.NewKeyValues(downloadsDir, kevlar_legacy.JsonExt)
	if err != nil {
		panic(err)
	}

	if _, err := kv.Cut(key); err != nil {
		panic(err)
	}

	if err := kv.Set(key, strings.NewReader("test")); err != nil {
		panic(err)
	}

	fmt.Println("kv.ModTime:", kv.ModTime())
	fmt.Println("kv.IsUpdatedAfter:", start, kv.IsUpdatedAfter(key, start))

	if err := kv.Set(key, strings.NewReader("this is not a test")); err != nil {
		panic(err)
	}

	fmt.Println("kv.ModTime:", kv.ModTime())
	fmt.Println("kv.IsUpdatedAfter:", start, kv.IsUpdatedAfter(key, start))

	fmt.Println("kv.Keys:", kv.Len())
	for id := range kv.Keys() {
		fmt.Println(" ", id)
	}

}
