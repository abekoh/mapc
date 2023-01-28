package main

import (
	"log"
	"reflect"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/examples/sqlc-domainmodel/domain"
	"github.com/abekoh/mapc/examples/sqlc-domainmodel/infrastructure/sqlc"
	"github.com/abekoh/mapc/mapcstd"
)

func main() {
	m := mapc.New()
	m.Option.OutPath = "infrastructure/mapper.go"
	// FIXME
	m.Option.TypeMappers = append(m.Option.TypeMappers, &mapcstd.MapTypeMapper{
		mapcstd.NewTyp(reflect.TypeOf("string")): map[*mapcstd.Typ]mapcstd.Caster{
			mapcstd.NewTyp(reflect.TypeOf(domain.UserID{})): mapcstd.NewSimpleCaster(&mapcstd.Caller{
				PkgPath:    "github.com/abekoh/mapc/examples/sqlc-domainmodel/domain",
				Name:       "UserID",
				CallerType: mapcstd.Type,
			}),
		},
	})

	m.Register(sqlc.User{}, domain.User{})
	m.Register(sqlc.Task{}, domain.Task{})
	m.Register(sqlc.SubTask{}, domain.SubTask{})

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
