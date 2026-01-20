#!/bin/sh

if [ ! -d bin/$(go env GOOS)_$(go env GOARCH) ]; then
   mkdir -p bin/$(go env GOOS)_$(go env GOARCH)
fi

go build -o bin/$(go env GOOS)_$(go env GOARCH)/convert ./cmd/convert $*
