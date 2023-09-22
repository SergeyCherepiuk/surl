package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type accountGetter struct {
	getStmt *sqlx.NamedStmt
}

func NewAccountGetter() *accountGetter {
	return &accountGetter{
		getStmt: internal.MustPrepare(db, `SELECT * FROM users WHERE username = :username`),
	}
}

func (ag accountGetter) Get(ctx context.Context, username string) (domain.User, error) {
	params := map[string]any{"username": username}
	user := domain.User{}
	err := ag.getStmt.SelectContext(ctx, &user, params)
	return user, err
}
