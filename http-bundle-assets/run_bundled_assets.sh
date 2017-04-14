#!/bin/bash
set -u -e -o pipefail
rm assets.zip || true
go run create_assets_zip/*.go
go run main.go -use-bundled-assets || true
