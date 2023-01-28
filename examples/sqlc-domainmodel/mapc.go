package main

import (
	"log"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/examples/sqlc-domainmodel/domain"
	"github.com/abekoh/mapc/examples/sqlc-domainmodel/infrastructure/sqlc"
)

func main() {
	m := mapc.New()
	m.Option.OutPath = "infrastructure/mapper.go"
	m.Register(sqlc.User{}, domain.User{})
	gs, errs := m.Generate()
	if len(errs) > 0 {
		log.Fatal(errs)
	}
	for _, g := range gs {
		err := g.Save()
		if err != nil {
			log.Fatal()
		}
	}
}
