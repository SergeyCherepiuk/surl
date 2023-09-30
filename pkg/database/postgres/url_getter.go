package postgres

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type urlGetter struct {
	getStmt                  *sqlx.NamedStmt
	getAllStmt               *sqlx.NamedStmt
	getAllSortedStmt         map[string]*sqlx.NamedStmt
	getAllSortedReversedStmt map[string]*sqlx.NamedStmt
}

func NewUrlGetter() *urlGetter {
	return &urlGetter{
		getStmt:                  internal.MustPrepare(db, `SELECT * FROM urls WHERE username = :username AND hash = :hash`),
		getAllStmt:               internal.MustPrepare(db, `SELECT * FROM urls WHERE username = :username`),
		getAllSortedStmt:         internal.MustPrepareMap(db, []string{"origin", "hash", "created_at", "last_used_at", "expires_at"}, `SELECT * FROM urls WHERE username = :username ORDER BY %s`),
		getAllSortedReversedStmt: internal.MustPrepareMap(db, []string{"origin", "hash", "created_at", "last_used_at", "expires_at"}, `SELECT * FROM urls WHERE username = :username ORDER BY %s DESC`),
	}
}

func (ug urlGetter) Get(ctx context.Context, username, hash string) (domain.Url, error) {
	params := map[string]any{
		"username": username,
		"hash":     hash,
	}
	var url domain.Url
	err := ug.getStmt.GetContext(ctx, &url, params)
	return url, err
}

func (ug urlGetter) GetAll(ctx context.Context, username string) ([]domain.Url, error) {
	params := map[string]any{"username": username}
	urls := []domain.Url{}
	err := ug.getAllStmt.SelectContext(ctx, &urls, params)
	return urls, err
}

func (ug urlGetter) GetAllSorted(ctx context.Context, username, sortBy string, reversed bool) ([]domain.Url, error) {
	params := map[string]any{"username": username}

	var stmt *sqlx.NamedStmt
	var ok bool
	if reversed {
		stmt, ok = ug.getAllSortedReversedStmt[sortBy]
	} else {
		stmt, ok = ug.getAllSortedStmt[sortBy]
	}

	if !ok {
		return []domain.Url{}, fmt.Errorf("sorting by given attribute is unsupported")
	}

	urls := []domain.Url{}
	err := stmt.SelectContext(ctx, &urls, params)
	return urls, err
}
