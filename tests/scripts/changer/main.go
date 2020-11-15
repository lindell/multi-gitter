package main

import (
	"bytes"
	"io/ioutil"
)

const fileName = "test.txt"

func main() {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	replaced := bytes.ReplaceAll(data, []byte("apple"), []byte("banana"))

	err = ioutil.WriteFile(fileName, replaced, 0600)
	if err != nil {
		panic(err)
	}
}
