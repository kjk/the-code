#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build
$now = Get-Date -UFormat "%Y-%m-%d_%T"
$sha1 = (git rev-parse HEAD).Trim()

go build -ldflags "-X main.sha1ver=$sha1 -X main.buildTime=$now"

Write-Host "Run:"
Write-Host ".\embed-build-number.exe -version : print version to stdout and exit"
Write-Host ".\embed-build-number.exe : start as web server and provide /app/debug page"
