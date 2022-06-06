#!/bin/bash
GOARCH=wasm GOOS=js go build -o igop.wasm
gopherjs build -v -m -o igop.js
