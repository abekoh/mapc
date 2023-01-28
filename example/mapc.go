package main

import (
	"log"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/example/domain"
	"github.com/abekoh/mapc/example/infrastructure/sqlc"
)

func main() {
	m := mapc.New()
	m.Option.OutPath = "infrastructure/mapper.go"
	m.Option.Mode = mapc.Deterministic

	// TODO: set UserID like `type UserID uuid.UUID`, and setup custom caster
	//m.Option.TypeMappers = append(m.Option.TypeMappers, mapcstd.TypeMapperFunc(
	//	func(src, dest *mapcstd.Typ) (mapcstd.Caster, bool) {
	//		if src.Equals(mapcstd.NewTypOf("")) && dest.Equals(mapcstd.NewTypOf(domain.UserID(""))) {
	//			return mapcstd.NewSimpleCaster(&mapcstd.Caller{
	//				PkgPath:    "github.com/abekoh/mapc/example/sqlc-domainmodel/domain",
	//				Name:       "UserID",
	//				CallerType: mapcstd.Type,
	//			}), true
	//		}
	//		return nil, false
	//	}),
	//)

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
