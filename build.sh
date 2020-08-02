#!/bin/bash
GOARCH=wasm GOOS=js go build -o igop.wasm
gopherjs build -m -o igop.js
