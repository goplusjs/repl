# repl.xgo.dev

This implements a web viewable version of the XGo REPL.

This is done by compiling gpython into wasm and running that in the
browser.

<https://goplusjs.github.io/repl/>


## Build and run

Build GoPlus REPL for GopherJS/WASM
```
go get github.com/goplusjs/gopherjs
git clone https://github.com/goplusjs/repl
cd repl
./build.sh
```


## Thanks

Thanks to [jQuery Terminal](https://terminal.jcubic.pl/) for the
terminal emulator and the go team for great [wasm
support](https://github.com/golang/go/wiki/WebAssembly).

Thanks to [Gpython](https://github.com/go-python/gpython) for REPL cli/web code.
