/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Package repl implements the ``gop repl'' command.
package main

import (
	"fmt"
	"io"

	"github.com/goplusjs/repl"
	"github.com/peterh/liner"
)

type LinerUI struct {
	state  *liner.State
	prompt string
}

func (u *LinerUI) SetPrompt(prompt string) {
	u.prompt = prompt
}

func (u *LinerUI) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

const (
	welcome string = "welcome to Go+ console!"
)

func main() {
	fmt.Println(welcome)

	state := liner.NewLiner()
	defer state.Close()
	state.SetCtrlCAborts(true)
	state.SetMultiLineMode(true)

	ui := &LinerUI{state: liner.NewLiner()}
	repl := repl.New(ui)
	for {
		line, err := ui.state.Prompt(ui.prompt)
		if err != nil {
			if err == liner.ErrPromptAborted || err == io.EOF {
				fmt.Printf("\n")
				break
			}
			fmt.Printf("Problem reading line: %v\n", err)
			continue
		}
		repl.Run(line)
	}
}
