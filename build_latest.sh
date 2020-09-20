#!/bin/bash
go mod edit -require=github.com/goplus/gop@latest
go list --tags js
GOARCH=wasm GOOS=js go build -o igop.wasm
gopherjs build -m -o igop.js
