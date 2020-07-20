#!/bin/bash

# Title: Fix linting problems in all your go repositories

golangci-lint run ./... --fix
