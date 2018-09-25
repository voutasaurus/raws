package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()

	if err := r(); err != nil {
		log.Fatal(err)
	}
}

func r() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("raws: error getting working directory: %v", err)
	}
	if err := filepath.Walk(dir, raw(dir)); err != nil {
		return fmt.Errorf("raws: error walking filesystem %v", err)
	}
	return nil
}

func raw(wd string) filepath.WalkFunc {
	f := func(path string, info os.FileInfo, err error) error {
		n := info.Name()
		ignore := n == "" || n[0] == '.' || n[0] == '_' || n == "vendor"

		if info.IsDir() {
			if ignore {
				return filepath.SkipDir
			}
			return nil
		}
		if !ignore && n[len(n)-3:] == ".go" {
			return rawf(wd, path)
		}
		return nil
	}
	return filepath.WalkFunc(f)
}

const (
	escape = "\x1b"

	setPink  = escape + "[35m"
	setGreen = escape + "[32m"
	reset    = escape + "[0m"
)

func rawf(wd, path string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}

	var out strings.Builder
	out.WriteString(setPink)
	out.WriteString(path[len(wd)+1:])
	out.WriteString("\n")
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BasicLit:
			if x.Kind == token.STRING && x.Value[0] == '`' {
				found = true
				pos := fset.Position(x.ValuePos)
				out.WriteString(setGreen)
				out.WriteString(fmt.Sprintf("%4d", pos.Line))
				out.WriteString(reset)
				out.WriteString(fmt.Sprintf(": %s\n", x.Value))
			}
		}
		return true
	})
	if found {
		out.WriteString("\n")
		os.Stdout.WriteString(out.String())
	}
	return nil
}
