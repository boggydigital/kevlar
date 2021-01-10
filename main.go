package main

import (
	"fmt"
	"github.com/boggydigital/kvas/internal"
	"time"
)

func main() {
	vsProducts, err := internal.NewJsonClient("products")
	if err != nil {
		fmt.Println(err.Error())
	}
	if vsProducts == nil {
		fmt.Println("couldn't load values set")
		return
	}

	key := "1"

	err = vsProducts.Set(key, []byte(key))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if vsProducts.Contains(key) {
		bytes, err := vsProducts.Get(key)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("value by the key '%s':%s\n", key, string(bytes))
	} else {
		fmt.Printf("products don't contain value by the key '%s'\n", key)
	}

	//if err := vsProducts.Remove(key); err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	fmt.Println("all:", vsProducts.CreatedAfter(time.Now().Unix()))
}
