#!/bin/bash

# Title: Replace all instances of empty interface with any

gofmt -r 'interface{} -> any' -w **/*.go
