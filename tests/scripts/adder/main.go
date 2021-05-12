package main

import (
	"flag"
	"io/ioutil"
	"strings"
)

func main() {
	filenames := flag.String("filenames", "", "")
	data := flag.String("data", "", "")
	flag.Parse()

	if *filenames == "" {
		panic("empty filename")
	}
	if *data == "" {
		panic("empty data")
	}

	for _, fn := range strings.Split(*filenames, ",") {
		err := ioutil.WriteFile(fn, []byte(*data), 0600)
		if err != nil {
			panic(err)
		}
	}
}
