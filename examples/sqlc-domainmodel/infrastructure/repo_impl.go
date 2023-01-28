package infrastructure

import (
	"context"
	"database/sql"

	"github.com/abekoh/mapc/examples/sqlc-domainmodel/domain"
	"github.com/abekoh/mapc/examples/sqlc-domainmodel/infrastructure/sqlc"
)

type RepositoryImpl struct {
	db *sql.DB
}

func (r RepositoryImpl) Get(ctx context.Context, id domain.UserID) (*domain.User, error) {
	queries := sqlc.New(r.db)
	user, err := queries.GetUser(ctx, id.String())
	if err != nil {
		return nil, err
	}
	res := MapUserToUser(user)
	return &res, nil
}

func (r RepositoryImpl) List(ctx context.Context) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}
