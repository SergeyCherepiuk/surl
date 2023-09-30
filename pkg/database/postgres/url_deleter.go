package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type urlDeleter struct {
	deleteStmt *sqlx.NamedStmt
}

func NewUrlDeleter() *urlDeleter {
	return &urlDeleter{
		deleteStmt: internal.MustPrepare(db, `DELETE FROM urls WHERE username = :username AND hash = :hash`),
	}
}

func (ud urlDeleter) Delete(ctx context.Context, username, hash string) error {
	params := map[string]any{
		"username": username,
		"hash":     hash,
	}
	_, err := ud.deleteStmt.ExecContext(ctx, params)
	return err
}
