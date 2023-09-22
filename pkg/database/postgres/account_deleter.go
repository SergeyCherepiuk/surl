package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type accountDeleter struct {
	deleteStmt *sqlx.NamedStmt
}

func NewAccountDeleter() *accountDeleter {
	return &accountDeleter{
		deleteStmt: internal.MustPrepare(db, `DELETE FROM users WHERE username = :username`),
	}
}

func (ad accountDeleter) Delete(ctx context.Context, username string) error {
	params := map[string]any{
		"username": username,
	}
	_, err := ad.deleteStmt.ExecContext(ctx, params)
	return err
}
