package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type verificator struct {
	verifyStmt *sqlx.NamedStmt
}

func NewVerificator() *verificator {
	return &verificator{
		verifyStmt: internal.MustPrepare(db, `UPDATE users SET is_verified = true WHERE username = :username`),
	}
}

func (v verificator) Verify(ctx context.Context, username string) error {
	params := map[string]any{"username": username}
	_, err := v.verifyStmt.ExecContext(ctx, params)
	return err
}
