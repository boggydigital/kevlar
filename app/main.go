package main

import (
	"fmt"
	"github.com/boggydigital/kvas"
	"time"
)

func main() {
	vs, err := kvas.NewClient("test", ".json")
	if err != nil {
		panic(err.Error())
	}
	if vs == nil {
		panic("couldn't load values set")
	}

	key := "value1"
	ts := time.Now().Unix()

	err = vs.Set(key, []byte(key))
	if err != nil {
		panic(err.Error())
	}

	if vs.Contains(key) {
		bytes, err := vs.Get(key)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("value by the key '%s':%s\n", key, string(bytes))
	} else {
		fmt.Printf("products don't contain value by the key '%s'\n", key)
	}

	fmt.Println("all:", vs.ModifiedAfter(ts))
}
