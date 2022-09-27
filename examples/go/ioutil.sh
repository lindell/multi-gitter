#!/bin/bash

# Title: Fix the ioutil deprecation

gofmt -w -r 'ioutil.Discard -> io.Discard' .
gofmt -w -r 'ioutil.NopCloser -> io.NopCloser' .
gofmt -w -r 'ioutil.ReadAll -> io.ReadAll' .
gofmt -w -r 'ioutil.ReadFile -> os.ReadFile' .
gofmt -w -r 'ioutil.TempDir -> os.MkdirTemp' .
gofmt -w -r 'ioutil.TempFile -> os.CreateTemp' .
gofmt -w -r 'ioutil.WriteFile -> os.WriteFile' .
gofmt -w -r 'ioutil.ReadDir -> os.ReadDir ' . # (note: returns a slice of os.DirEntry rather than a slice of fs.FileInfo)

goimports -w .
