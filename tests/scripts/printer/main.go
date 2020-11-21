package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const fileName = "test.txt"

func main() {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	fmt.Fprintln(os.Stderr, strings.ToUpper(string(data)))
}
