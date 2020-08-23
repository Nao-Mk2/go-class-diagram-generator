package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Nao-Mk2/go-class-diagram-generator/internal/generator"
	"github.com/Nao-Mk2/go-class-diagram-generator/internal/gocdparser"
)

func main() {
	var recursion, includeStandard, includeTest bool
	flag.BoolVar(&recursion, "r", false, "")
	flag.BoolVar(&recursion, "recursive", false, "search sub-directory recursively")
	flag.BoolVar(&includeStandard, "s", false, "")
	flag.BoolVar(&includeStandard, "standard", false, "include standard library packages")
	flag.BoolVar(&includeTest, "t", false, "")
	flag.BoolVar(&includeTest, "test", false, "include _test packages")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println(`Usage:

	gocdgen <path>
	
The flags are:

	r, recursive search sub-directory recursively`)
		os.Exit(0)
	}

	path, err := getAbsolutePath()
	if err != nil {
		log.Fatalf("failure to get absolute path: %+v", err)
	}

	parser := gocdparser.GoCodeParser{
		IncludeTest: includeTest,
		Recursion:   recursion,
	}
	pkgs, err := parser.ParsePackages(path)
	if err != nil {
		log.Fatalf("failure to parse packages: %+v", err)
	}

	gen := generator.StdOutGenerator{
		Packages:        pkgs,
		IncludeStandard: includeStandard,
	}
	gen.Generate()
}

func getAbsolutePath() (string, error) {
	p := flag.Args()[0]
	fi, err := os.Stat(p)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("could not find directory or file: %s", p)
	}
	if !fi.IsDir() && !fi.Mode().IsRegular() {
		return "", fmt.Errorf("not a directory or file: %s", p)
	}

	abs, e := filepath.Abs(p)

	return abs, e
}
