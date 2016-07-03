package mold

import (
	"go/parser"
	"go/token"
	"io"
	"os"
)

// LoadTypes loads all top-level functions and symbols from a source file
func LoadTypes(r io.Reader) (map[string]Type, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "src.go", r, 0)
	if err != nil {
		return nil, err
	}

	b := newBuilder(file.Name.Name)
	for _, decl := range file.Decls {
		b.add(decl)
	}
	b.build()
	return b.named, nil
}

// LoadFile loads all top-level functions and symbols from a source file
func LoadFile(path string) (map[string]Type, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, f, 0)
	if err != nil {
		return nil, err
	}

	b := newBuilder(file.Name.Name)
	for _, decl := range file.Decls {
		b.add(decl)
	}
	b.build()
	return b.named, nil
}
