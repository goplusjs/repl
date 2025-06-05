// An online REPL for gpython using wasm

// Copyright 2018 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js
// +build js

package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"syscall/js"

	"github.com/goplus/ixgo"
	"github.com/goplus/ixgo/repl"
	_ "github.com/goplus/ixgo/xgobuild"
)

// Implement the replUI interface
type termIO struct {
	js.Value
}

// SetPrompt sets the UI prompt
func (t *termIO) SetPrompt(prompt string) {
	t.Call("set_prompt", prompt)
}

// Print outputs the string to the output
func (t *termIO) Printf(format string, a ...interface{}) {
	line := fmt.Sprintf(format, a...)
	t.Call("echo", strings.TrimRight(line, "\n"))
}

var document js.Value

func getElementById(name string) js.Value {
	node := document.Call("getElementById", name)
	if node.IsUndefined() {
		log.Fatalf("Couldn't find element %q\n", name)
	}
	return node
}

func running() string {
	switch {
	case runtime.GOOS == "js" && runtime.GOARCH == "wasm":
		return "Wasm"
	case runtime.Compiler == "gopherjs":
		return "GopherJS"
	}
	return "unknown"
}

const (
	isGopherJS = runtime.GOARCH == "ecmascript"
)

func main() {
	document = js.Global().Get("document")
	if document.IsUndefined() {
		log.Fatalf("Didn't find document - not running in browser\n")
	}

	// Clear the loading text
	termNode := getElementById("term")
	termNode.Set("innerHTML", "")

	// work out what we are running on and mark active
	tech := running()
	gopVersion := getElementById("GopVersion").Get("innerHTML").String()
	iGopVersion := getElementById("iGopVersion").Get("innerHTML").String()
	node := getElementById(tech)
	node.Get("classList").Call("add", "active")

	// Make a repl referring to an empty term for the moment
	REPL := repl.NewREPL(0)
	REPL.SetFileName("main.gop")

	var term *termIO
	gopCheck := js.Global().Call("$", "#enableGoplus")

	getGreetings := func(gop bool) string {
		var mode string
		if gop {
			mode = "XGo"
		} else {
			mode = "Go"
		}
		return fmt.Sprintf("iXGo %v (%v, gop %v) running in your browser with %v. (%v Mode)",
			iGopVersion, runtime.Version(), gopVersion, tech, mode)
	}

	gopCheck.Call("change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		REPL.Repl = ixgo.NewRepl(ixgo.NewContext(0))
		gop := this.Get("checked").Bool()
		if gop {
			REPL.SetFileName("main.gop")
		} else {
			REPL.SetFileName("main.go")
		}
		term.Call("clear")
		term.Call("echo", getGreetings(gop))
		return nil
	}))

	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err := REPL.Run(args[0].String())
		if err != nil {
			term.Printf("%v\n", err)
		}
		return nil
	})

	// Create a jquery terminal instance
	opts := js.ValueOf(map[string]interface{}{
		"greetings": getGreetings(true),
		"name":      "goplus",
		"prompt":    repl.NormalPrompt,
		"clear":     true,
	})
	terminal := js.Global().Call("$", "#term").Call("terminal", cb, opts)
	term = &termIO{terminal}
	// Send the console log direct to the terminal
	js.Global().Get("console").Set("log", terminal.Get("echo"))

	// Set the implementation of term
	REPL.SetUI(term)

	// wait for callbacks
	select {}
}

var (
	GopVersion  = "v1.1.2"
	iGopVersion = "v0.9.4"
)
