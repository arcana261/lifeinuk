#!/bin/bash

cat data/highlights.txt | tr '\n' ' ' | sed 's|---|\n|g' | fzf --preview 'echo {} | fold -w 40 -s'
