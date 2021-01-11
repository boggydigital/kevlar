package main

import (
	"fmt"
	"github.com/boggydigital/kvas"
)

func main() {
	vs, err := kvas.NewClient("test", ".json", true)
	if err != nil {
		panic(err.Error())
	}
	if vs == nil {
		panic("couldn't load values set")
	}

	key := "value1"

	err = vs.Set(key, []byte(key))
	if err != nil {
		panic(err.Error())
	}

	bytes, err := vs.Get(key)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("value by the key '%s':%s\n", key, string(bytes))

	fmt.Println("all:", vs.All())
}
