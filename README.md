# repl.goplus.org

This implements a web viewable version of the GoPlus REPL.

This is done by compiling gpython into wasm and running that in the
browser.

https://goplusjs.github.io/repl/


## Build and run

Run GoPlus REPL for GopherJS/WASM.

`make build` will build with go GopherJS/WASM (you'll need go1.14 minimum)

`make serve` will build and run a local webserver you can see the results on.


## Thanks

Thanks to [jQuery Terminal](https://terminal.jcubic.pl/) for the
terminal emulator and the go team for great [wasm
support](https://github.com/golang/go/wiki/WebAssembly).

Thanks to [Gpython](https://github.com/go-python/gpython) for REPL cli/web code.
