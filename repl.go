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
package repl

import (
	"errors"
	"strings"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"

	exec "github.com/goplus/gop/exec/bytecode"
	_ "github.com/goplus/gop/lib"
)

type UI interface {
	SetPrompt(prompt string)
	Printf(format string, a ...interface{})
}

type REPL struct {
	src            string       // the whole source code from repl
	preContext     exec.Context // store the context after exec
	ip             int          // store the ip after exec
	normalPrompt   string       // the prompt type in console
	continuePrompt string
	continueMode   bool // switch to control the promot type
	term           UI   // liner instance
}

const (
	continuePrompt string = "... "
	standardPrompt string = ">>> "
)

func New(term UI) *REPL {
	r := &REPL{normalPrompt: standardPrompt, continuePrompt: continuePrompt, term: term}
	term.SetPrompt(r.normalPrompt)
	return r
}

// Run handle one line
func (r *REPL) Run(line string) {
	if r.continueMode {
		r.continueModeByLine(line)
	}
	if !r.continueMode {
		r.run(line)
	}
}

// continueModeByLine check if continue-mode should continue :)
func (r *REPL) continueModeByLine(line string) {
	if line != "" {
		r.src += line + "\n"
		return
	}
	// input nothing means jump out continue mode
	r.continueMode = false
	r.term.SetPrompt(r.normalPrompt)
}

// run execute the input line
func (r *REPL) run(newLine string) (err error) {
	src := r.src + newLine + "\n"
	defer func() {
		if err == nil {
			r.src = src
		}
		if errR := recover(); errR != nil {
			r.dumpErr(newLine, errR)
			err = errors.New("panic err")
			// TODO: Need a better way to log and show the stack when crash
			// It is too long to print stack on terminal even only print part of the them; not friendly to user
		}
	}()
	fset := token.NewFileSet()
	pkgs, err := parser.Parse(fset, "", src, 0)
	if err != nil {
		// check if into continue mode
		if strings.Contains(err.Error(), `expected ')', found 'EOF'`) ||
			strings.Contains(err.Error(), "expected '}', found 'EOF'") {
			r.term.SetPrompt(r.continuePrompt)
			r.continueMode = true
			err = nil
			return
		}
		r.term.Printf("ParseGopFiles err: %v\n", err)
		return
	}
	cl.CallBuiltinOp = exec.CallBuiltinOp

	b := exec.NewBuilder(nil)

	_, err = cl.NewPackage(b.Interface(), pkgs["main"], fset, cl.PkgActClMain)
	if err != nil {
		if err == cl.ErrMainFuncNotFound {
			err = nil
			return
		}
		r.term.Printf("NewPackage err %+v\n", err)
		return
	}
	code := b.Resolve()
	ctx := exec.NewContext(code)
	if r.ip != 0 {
		// if it is not the first time, restore pre var
		r.preContext.CloneSetVarScope(ctx)
	}
	currentIP := ctx.Exec(r.ip, code.Len())
	r.preContext = *ctx
	// "currentip - 1" is the index of `return`
	// next time it will replace by new code from newLine
	r.ip = currentIP - 1
	return
}

func (r *REPL) dumpErr(line string, err interface{}) {
	r.term.Printf("code run fail : %v\n", line)
	r.term.Printf("repl err: %v\n", err)
}
