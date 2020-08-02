@echo on
set GOOS=darwin
gopherjs build -m -o igop.js
set GOARCH=wasm
set GOOS=js
go build -v -o igop.wasm
