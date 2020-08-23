package gocdparser

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nao-Mk2/go-class-diagram-generator/internal/entity"
)

// GoCodeParser implements parser methods for Go source files.
type GoCodeParser struct {
	IncludeTest bool
	Recursion   bool
}

// ParsePackages returns a map of package name -> package AST with all the packages found.
func (p GoCodeParser) ParsePackages(path string) ([]*entity.Package, error) {
	fi, e := os.Stat(path)
	if os.IsNotExist(e) {
		return nil, fmt.Errorf("could not find directory or file: %s", path)
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		paths, err := getDirs(path, p.Recursion)
		if err != nil {
			return nil, err
		}
		return parseDir(paths, p.IncludeTest)

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

func contains(slice []string, v string) bool {
	for _, e := range slice {
		if e == v {
			return true
		}
	}
	return false
}

func getDirs(path string, recursion bool) ([]string, error) {
	if !recursion {
		return []string{path}, nil
	}

	paths := make([]string, 0)
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("failure to accesse a path %q: %v\n", p, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		// append directory containing files
		dir := strings.Replace(p, info.Name(), "", 1)
		if !contains(paths, dir) {
			paths = append(paths, dir)
		}

		return nil
	})

	return paths, nil
}

func makeUnique(imps []entity.Import) []entity.Import {
	keys := make(map[string]bool, 0)
	unique := make([]entity.Import, 0)
	for _, imp := range imps {
		if _, v := keys[imp.Path]; !v {
			keys[imp.Path] = true
			unique = append(unique, imp)
		}
	}

	return unique
}

func parseDir(paths []string, includeTest bool) ([]*entity.Package, error) {

	pkgMap := make(map[string]*entity.Package, 0)
	for _, path := range paths {
		fs := token.NewFileSet()
		parsed, err := parser.ParseDir(fs, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil, err
		}

		for _, p := range parsed {
			if strings.Index(p.Name, "_test") > -1 && !includeTest {
				continue
			}

			for _, f := range p.Files {
				pkg, err := entity.NewPackage(f.Name.Name, f.Imports)
				if err != nil {
					return nil, err
				}

				if pkgMap[f.Name.Name] != nil {
					pkgMap[f.Name.Name].Imports = append(pkgMap[f.Name.Name].Imports, pkg.Imports...)
				} else {
					pkgMap[f.Name.Name] = pkg
				}
			}
		}
	}

	pkgs := make([]*entity.Package, 0)
	for _, v := range pkgMap {
		pkgs = append(pkgs, v)
	}

	for _, pkg := range pkgs {
		pkg.Imports = makeUnique(pkg.Imports)
	}

	return pkgs, nil
}

func parseFile(path string) (*entity.Package, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	return entity.NewPackage(f.Name.Name, f.Imports)
}
