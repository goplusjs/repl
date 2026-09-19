#!/bin/bash
rm docs/*.js
rm docs/*.js.map
rm docs/*.wasm
cp $GOROOT/lib/wasm/wasm_exec.js docs/
GOFLAGS='-tags=goplus.ifacefuncval' go run make.go
