package main

import (
	"fmt"
	"io/ioutil"
)

const fileName = "test.txt"

func main() {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
