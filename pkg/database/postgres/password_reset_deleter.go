package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type passwordResetDeleter struct {
	deleteAllStmt *sqlx.NamedStmt
}

func NewPasswordResetDeleter() *passwordResetDeleter {
	return &passwordResetDeleter{
		deleteAllStmt: internal.MustPrepare(db, `DELETE FROM password_reset_requests WHERE username = :username`),
	}
}

func (prd passwordResetDeleter) DeleteAll(ctx context.Context, username string) error {
	params := map[string]any{"username": username}
	_, err := prd.deleteAllStmt.ExecContext(ctx, params)
	return err
}
