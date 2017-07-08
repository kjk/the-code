#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build
now=$(date +'%Y-%m-%d_%T')
go build -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"

echo "Run:"
echo "./embed-build-number -version : print version to stdout and exit"
echo "./embed-build-number : start as web server and provide /app/debug page"
