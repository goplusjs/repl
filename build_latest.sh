#!/bin/bash
go mod edit -require=github.com/goplus/gop@latest
go mod download github.com/goplus/gop
GOARCH=wasm GOOS=js go build -o igop.wasm
gopherjs build -v -m -o igop.js
