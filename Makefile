build:
	GOARCH=wasm GOOS=js go build -o igop.wasm
	gopherjs build -m -o igop.js

serve:	build
	go run serve.go
