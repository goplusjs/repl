//go:build ingore
// +build ingore

package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type Module struct {
	Path      string
	Version   string
	Time      time.Time
	Dir       string
	GoMod     string
	GoVersion string
}

func getModule(path string) (*Module, error) {
	cmd := exec.Command("go", "list", "-m", "-json", path)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var m Module
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if m.Version == "" {
		m.Version = getLocalModuleVersion(m.Dir)
	}
	return &m, err
}

func main() {
	xgo, err := getModule("github.com/goplus/xgo")
	check(err)
	ixgo, err := getModule("github.com/goplus/ixgo")
	check(err)

	tag, err := getHash()
	fmt.Println(tag, xgo.Version, ixgo.Version)

	if err != nil {
		panic(err)
	}
	// build index
	data, err := ioutil.ReadFile("./index_tpl.html")
	check(err)
	data = bytes.Replace(data, []byte("loader.js"), []byte("loader_"+tag+".js"), 1)
	data = bytes.Replace(data, []byte("$XgoVersion"), []byte(xgo.Version), 1)
	data = bytes.Replace(data, []byte("$ixgoVersion"), []byte(ixgo.Version), 1)
	data = bytes.Replace(data, []byte("$GoVersion"), []byte("Go "+goVersion()), 1)
	err = ioutil.WriteFile("./docs/index.html", data, 0755)

	// build loader.js
	data, err = ioutil.ReadFile("./loader_tpl.js")
	check(err)

	data = bytes.Replace(data, []byte("igop"), []byte("ixgo_"+tag), 2)
	err = ioutil.WriteFile("./docs/loader_"+tag+".js", data, 0755)
	check(err)

	// err = build_js("./docs", "ixgo_"+tag)
	// check(err)

	err = build_wasm("./docs", "ixgo_"+tag)
	check(err)
}

func goVersion() string {
	version := runtime.Version()
	if len(version) > 2 && version[:2] == "go" {
		return version[2:]
	}
	return version
}

func getLocalModuleVersion(dir string) string {
	if dir != "" {
		cmd := exec.Command("git", "-C", dir, "describe", "--tags", "--always", "--dirty")
		if data, err := cmd.Output(); err == nil {
			if version := string(bytes.TrimSpace(data)); version != "" {
				return version
			}
		}
	}
	return "(devel)"
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func getHash() (string, error) {
	h := md5.New()
	for _, f := range []string{"main.go", "pkg_std.go", "pkg_gop.go", "go.mod"} {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return "", err
		}
		h.Write(data)
	}
	return fmt.Sprintf("%x", h.Sum(nil)[:4]), nil
	// cmd := exec.Command("git", "describe", "--tag")
	// return cmd.Output()
}

// GOARCH=wasm GOOS=js go build -o igop.wasm
// gopherjs build -v -m -o igop.js

func build_js(dir, tag string) error {
	cmd := exec.Command("gopherjs", "build", "-v", "-m", "-o", filepath.Join(dir, tag+".js"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env := os.Environ()
	cmd.Env = append(env, "GOARCH=ecmascript", "GOOS=js")
	return cmd.Run()
}

func build_wasm(dir, tag string) error {
	//cmd := exec.Command("go", "build", "-ldflags", "-checklinkname=0", "-o", filepath.Join(dir, tag+".wasm"))
	cmd := exec.Command("go", "build", "-o", filepath.Join(dir, tag+".wasm"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env := os.Environ()
	cmd.Env = append(env, "GOARCH=wasm", "GOOS=js")
	return cmd.Run()
}
