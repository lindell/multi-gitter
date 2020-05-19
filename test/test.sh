#!/bin/bash

FILE="README.md"
if [ ! -f "$FILE" ]; then
    exit 1
fi

echo "Some extra test" >> README.md
