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
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println(`Usage:

	gocdgen <path>`)
		os.Exit(0)
	}

	path, err := getAbsolutePath()
	if err != nil {
		log.Fatalf("failure to get absolute path: %+v", err)
	}

	pkgs, err := gocdparser.ParsePackages(path)
	if err != nil {
		log.Fatalf("failure to parse packages: %+v", err)
	}

	gen := generator.StdOutGenerator{
		Packages: pkgs,
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
