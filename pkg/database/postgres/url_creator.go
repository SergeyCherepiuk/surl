package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type urlCreator struct {
	createStmt *sqlx.NamedStmt
}

func NewUrlCreator() *urlCreator {
	return &urlCreator{
		createStmt: internal.MustPrepare(db, `INSERT INTO urls VALUES (:username, :hash, :origin, :expires_at)`),
	}
}

func (uc urlCreator) Create(ctx context.Context, url domain.Url) error {
	params := map[string]any{
		"username":   url.Username,
		"hash":       url.Hash,
		"origin":     url.Origin,
		"expires_at": url.ExpiresAt,
	}
	_, err := uc.createStmt.ExecContext(ctx, params)
	return err
}
