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

	gojs "github.com/gopherjs/gopherjs/js"
	_ "github.com/goplus/igop/gopbuild"
	"github.com/goplus/igop/repl"
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
	if isGopherJS {
		gojs.Global.Call("$flushConsole")
	}
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
	case runtime.GOARCH == "ecmascript":
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
	node := getElementById(tech)
	node.Get("classList").Call("add", "active")

	// Make a repl referring to an empty term for the moment
	REPL := repl.NewREPL(0)
	REPL.SetFileName("main.gop")

	var term *termIO

	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err := REPL.Run(args[0].String())
		if err != nil {
			term.Printf("%v\n", err)
		}
		return nil
	})

	gopVersion := getElementById("GopVersion").Get("innerHTML").String()
	iGopVersion := getElementById("iGopVersion").Get("innerHTML").String()

	// Create a jquery terminal instance
	opts := js.ValueOf(map[string]interface{}{
		"greetings": fmt.Sprintf("iGo+ %v (gop version %v) running in your browser with %v",
			iGopVersion, gopVersion, tech),
		"name":   "goplus",
		"prompt": repl.NormalPrompt,
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
