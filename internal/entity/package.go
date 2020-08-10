package entity

import (
	"errors"
	"go/ast"
	"strings"
)

// Package represents a Go package.
type (
	Package struct {
		Name    string
		Imports []Import
	}

	// Import represents a imported Go library.
	Import struct {
		Path              string
		IsStandardLibrary bool
	}
)

// NewPackage creates a new Package.
func NewPackage(name string, imports []*ast.ImportSpec) (*Package, error) {
	if name == "" {
		return nil, errors.New("package name must not be empty")
	}

	imps := make([]Import, 0, len(imports))
	for _, imp := range imports {
		imps = append(imps, newImport(imp.Path.Value))
	}

	return &Package{
		Name:    name,
		Imports: imps,
	}, nil
}

func newImport(path string) Import {
	return Import{
		Path:              path,
		IsStandardLibrary: isStandardImportPath(path),
	}
}

// https://github.com/golang/go/blob/master/src/cmd/go/internal/search/search.go#L552
func isStandardImportPath(p string) bool {
	i := strings.Index(p, "/")
	if i < 0 {
		i = len(p)
	}

	return !strings.Contains(p[:i], ".")
}
