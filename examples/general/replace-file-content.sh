#!/bin/bash

# Title: Replace a file if it exist

REPLACE_FILE=~/test/pull_request_template.md # The file that should replace the file in the repo, must be an absolute path
FILE=.github/pull_request_template.md # Relative from any repos root

# Don't replace this file if it does not already exist in the repo
if [ ! -f "$FILE" ]; then
    exit 1
fi

cp $REPLACE_FILE $FILE
