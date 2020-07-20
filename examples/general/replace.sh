#!/bin/bash

# Title: Replace text in all files

# Assuming you are using gnu sed, if you are running this on a mac, please see https://stackoverflow.com/questions/4247068/sed-command-with-i-option-failing-on-mac-but-works-on-linux

find ./ -type f -exec sed -i -e 's/apple/orange/g' {} \;
