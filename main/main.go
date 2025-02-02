package main

import (
	"fmt"
	"github.com/boggydigital/kevlar"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {

	uhd, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	absDir := filepath.Join(uhd, "Downloads")

	kv, err := kevlar.New(absDir, kevlar.JsonExt)
	if err != nil {
		panic(err)
	}

	start := time.Now().UTC().Unix()

	for ii := range 10 {
		aa := strconv.Itoa(ii)
		if err := kv.Set(aa, strings.NewReader(aa)); err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second)

	for ii := range 5 {
		aa := strconv.Itoa(ii)
		if err := kv.Set(aa, strings.NewReader(aa+aa)); err != nil {
			panic(err)
		}
	}

	for ii := range 10 {
		aa := strconv.Itoa(ii)
		if err := kv.Cut(aa); err != nil {
			panic(err)
		}
	}

	if err := os.Remove(filepath.Join(absDir, "_log.gob")); err != nil {
		panic(err)
	}

	for id, mt := range kv.Since(start, kevlar.Create, kevlar.Update, kevlar.Cut) {
		fmt.Println(id, mt)
	}
}
