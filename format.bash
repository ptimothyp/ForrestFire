#!/usr/bin/env bash

for file in `git ls-files | egrep '.*\.go$'`; do
    gofmt -w $file
done
