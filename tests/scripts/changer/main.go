package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"time"
)

const fileName = "test.txt"

func main() {
	duration := flag.String("sleep", "", "Time to sleep before running the script")
	flag.Parse()

	if *duration != "" {
		d, _ := time.ParseDuration(*duration)
		time.Sleep(d)
	}

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
