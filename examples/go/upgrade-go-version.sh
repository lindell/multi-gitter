#!/bin/bash

# Title: Upgrade Go version in go modules

go mod edit -go 1.20
go mod tidy
