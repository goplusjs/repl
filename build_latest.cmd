@echo on
go mod edit -require=github.com/goplus/gop@latest
go list --tags js
gopherjs build -m -o igop.js
set GOARCH=wasm
set GOOS=js
go build -v -o igop.wasm
