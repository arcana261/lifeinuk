#!/bin/bash

cat scores.txt | cut -d ' ' -f 3 | sort -n | uniq -c

