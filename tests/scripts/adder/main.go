package main

import (
	"flag"
	"os"
	"path/filepath"
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
		dir := filepath.Dir(fn)
		if dir != "." {
			totalFilepath := "."
			for _, fp := range strings.Split(dir, string(filepath.Separator)) {
				totalFilepath = filepath.Join(totalFilepath, fp)
				err := os.Mkdir(totalFilepath, 0755)
				if err != nil {
					panic(err)
				}
			}
		}

		err := os.WriteFile(fn, []byte(*data), 0600)
		if err != nil {
			panic(err)
		}
	}
}
