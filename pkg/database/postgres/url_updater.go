package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type urlUpdater struct {
	updateStmt *sqlx.NamedStmt
}

func NewUrlUpdater() *urlUpdater {
	return &urlUpdater{
		updateStmt: internal.MustPrepare(db, `UPDATE urls SET origin = :new_origin, hash = :new_hash, last_used_at = :new_last_used_at WHERE username = :username AND hash = :hash`),
	}
}

func (uu urlUpdater) Update(ctx context.Context, username, hash string, updates domain.UrlUpdates) error {
	params := map[string]any{
		"username":         username,
		"hash":             hash,
		"new_origin":       updates.Origin,
		"new_hash":         updates.Hash,
		"new_last_used_at": updates.LastUsedAt,
	}
	_, err := uu.updateStmt.ExecContext(ctx, params)
	return err
}
