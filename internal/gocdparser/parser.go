package gocdparser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/Nao-Mk2/go-class-diagram-generator/internal/entity"
)

// ParsePackages returns a map of package name -> package AST with all the packages found.
func ParsePackages(path string) ([]*entity.Package, error) {
	fi, e := os.Stat(path)
	if os.IsNotExist(e) {
		return nil, fmt.Errorf("could not find directory or file: %s", path)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return parseDir(path)

	case mode.IsRegular():
		pkg, err := parseFile(path)
		if err != nil {
			return nil, err
		}
		return []*entity.Package{pkg}, nil

	default:
		return nil, fmt.Errorf("not a directory or file: %s", path)
	}
}

func parseFile(path string) (*entity.Package, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	return entity.NewPackage(f.Name.Name, f.Imports)
}

func parseDir(path string) ([]*entity.Package, error) {
	fs := token.NewFileSet()
	parsed, err := parser.ParseDir(fs, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	pkgs := make([]*entity.Package, 0)
	for _, p := range parsed {
		for _, f := range p.Files {
			pkg, err := entity.NewPackage(f.Name.Name, f.Imports)
			if err != nil {
				return nil, err
			}

			pkgs = append(pkgs, pkg)
		}
	}

	return pkgs, nil
}
