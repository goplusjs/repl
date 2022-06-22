#!/bin/bash
rm docs/*.js
rm docs/*.js.map
rm docs/*.wasm
cp wasm_exec.js docs/
go run make.go