#!/bin/bash
set -u -e -o pipefail
rm assets.zip || true
go run main.go
