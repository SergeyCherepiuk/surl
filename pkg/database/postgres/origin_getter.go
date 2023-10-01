package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type originGetter struct {
	getStmt *sqlx.NamedStmt
}

func NewOriginGetter() *originGetter {
	return &originGetter{
		getStmt: internal.MustPrepare(db, `SELECT * FROM urls WHERE username = :username AND hash = :hash`),
	}
}

func (og originGetter) Get(ctx context.Context, username, hash string) (string, time.Duration, error) {
	params := map[string]any{
		"username": username,
		"hash":     hash,
	}
	var url domain.Url

	if err := og.getStmt.GetContext(ctx, &url, params); err != nil {
		return "", time.Until(time.Now().In(time.UTC)), err
	}

	if url.ExpiresAt.Before(time.Now().In(time.UTC)) {
		return "", time.Until(time.Now().In(time.UTC)), fmt.Errorf("link expired")
	}

	return url.Origin, time.Until(url.ExpiresAt), nil
}
