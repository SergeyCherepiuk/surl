package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type accountCreator struct {
	createStmt *sqlx.NamedStmt
}

func NewAccountCreator() *accountCreator {
	return &accountCreator{
		createStmt: internal.MustPrepare(db, `INSERT INTO users VALUES (:username, :password)`),
	}
}

func (ac accountCreator) Create(ctx context.Context, user domain.User) error {
	params := map[string]any{
		"username": user.Username,
		"password": user.Password,
	}
	_, err := ac.createStmt.ExecContext(ctx, params)
	return err
}
