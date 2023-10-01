package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type verificationDeleter struct {
	deleteAllStmt *sqlx.NamedStmt
}

func NewVerificationDeleter() *verificationDeleter {
	return &verificationDeleter{
		deleteAllStmt: internal.MustPrepare(db, `DELETE FROM verification_requests WHERE username = :username`),
	}
}

func (vd verificationDeleter) DeleteAll(ctx context.Context, username string) error {
	params := map[string]any{"username": username}
	_, err := vd.deleteAllStmt.ExecContext(ctx, params)
	return err
}
