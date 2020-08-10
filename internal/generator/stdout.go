package generator

import (
	"log"

	"github.com/Nao-Mk2/go-class-diagram-generator/internal/entity"
)

type StdOutGenerator struct {
	Packages []*entity.Package
}

func (sg StdOutGenerator) Generate() {
	for _, pkg := range sg.Packages {
		log.Printf("- %s", pkg.Name)

		var stds, users []string
		for _, imp := range pkg.Imports {
			if imp.IsStandardLibrary {
				stds = append(stds, imp.Path)
			} else {
				users = append(users, imp.Path)
			}
		}

		if len(stds) > 0 {
			log.Println("--- standard libraries")
			for _, std := range stds {
				log.Printf("----- %s", std)
			}
		}
		if len(users) > 0 {
			log.Println("--- user libraries")
			for _, user := range users {
				log.Printf("----- %s", user)
			}
		}
	}
}
