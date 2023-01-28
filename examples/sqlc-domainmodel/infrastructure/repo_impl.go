package infrastructure

import (
	"context"
	"database/sql"

	"github.com/abekoh/mapc/examples/sqlc-domainmodel/domain"
)

type RepositoryImpl struct {
	db *sql.DB
}

func (r RepositoryImpl) Get(ctx context.Context, id domain.UserID) (*domain.User, error) {
	//queries := New(r.db)
	//TODO implement me
	panic("implement me")
}

func (r RepositoryImpl) List(ctx context.Context) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}
