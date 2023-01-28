package main

import (
	"log"

	"github.com/abekoh/mapc"
	"github.com/abekoh/mapc/example/domain"
	"github.com/abekoh/mapc/example/graph/model"
	"github.com/abekoh/mapc/example/infrastructure/sqlc"
)

func main() {
	m := mapc.New()

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

	infra := m.Group(func(option *mapc.Option) {
		option.OutPath = "infrastructure/mapper.go"
		option.Mode = mapc.Deterministic
	})

	infra.Register(sqlc.User{}, domain.User{})
	infra.Register(sqlc.Task{}, domain.Task{})
	infra.Register(sqlc.SubTask{}, domain.SubTask{})

	graph := m.Group(func(option *mapc.Option) {
		option.OutPath = "graph/mapper.go"
		option.Mode = mapc.Deterministic
	})
	graph.Register(sqlc.User{}, model.User{})
	graph.Register(sqlc.Task{}, model.Task{})
	graph.Register(sqlc.SubTask{}, model.SubTask{})

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
