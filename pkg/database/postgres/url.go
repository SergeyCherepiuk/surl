package postgres

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type urlService struct {
	getStmt                  *sqlx.NamedStmt
	getAllStmt               *sqlx.NamedStmt
	getAllSortedStmt         map[string]*sqlx.NamedStmt
	getAllSortedReversedStmt map[string]*sqlx.NamedStmt
	createStmt               *sqlx.NamedStmt
	updateStmt               *sqlx.NamedStmt
	deleteStmt               *sqlx.NamedStmt
}

func NewUrlService() *urlService {
	return &urlService{
		getStmt:                  internal.MustPrepare(db, `SELECT * FROM urls WHERE username = :username AND hash = :hash`),
		getAllStmt:               internal.MustPrepare(db, `SELECT * FROM urls WHERE username = :username`),
		getAllSortedStmt:         internal.MustPrepareMap(db, []string{"origin", "hash", "created_at", "last_used_at"}, `SELECT * FROM urls WHERE username = :username ORDER BY %s`),
		getAllSortedReversedStmt: internal.MustPrepareMap(db, []string{"origin", "hash", "created_at", "last_used_at"}, `SELECT * FROM urls WHERE username = :username ORDER BY %s DESC`),
		createStmt:               internal.MustPrepare(db, `INSERT INTO urls VALUES (:username, :hash, :origin)`),
		updateStmt:               internal.MustPrepare(db, `UPDATE urls SET origin = :origin, last_used_at = :last_used_at WHERE username = :username AND hash = :hash`),
		deleteStmt:               internal.MustPrepare(db, `DELETE FROM urls WHERE username = :username AND hash = :hash`),
	}
}

func (s urlService) Get(ctx context.Context, username, hash string) (domain.Url, error) {
	params := map[string]any{
		"username": username,
		"hash":     hash,
	}
	var url domain.Url
	err := s.getStmt.GetContext(ctx, &url, params)
	return url, err
}

func (s urlService) GetAll(ctx context.Context, username string) ([]domain.Url, error) {
	params := map[string]any{"username": username}
	urls := []domain.Url{}
	err := s.getAllStmt.SelectContext(ctx, &urls, params)
	return urls, err
}

func (s urlService) GetAllSorted(ctx context.Context, username, sortBy string, reversed bool) ([]domain.Url, error) {
	params := map[string]any{"username": username}

	var stmt *sqlx.NamedStmt
	var ok bool
	if reversed {
		stmt, ok = s.getAllSortedReversedStmt[sortBy]
	} else {
		stmt, ok = s.getAllSortedStmt[sortBy]
	}

	if !ok {
		return []domain.Url{}, fmt.Errorf("sorting by given attribute is unsupported")
	}

	urls := []domain.Url{}
	err := stmt.SelectContext(ctx, &urls, params)
	return urls, err
}

func (s urlService) Create(ctx context.Context, url domain.Url) error {
	params := map[string]any{
		"username": url.Username,
		"hash":     url.Hash,
		"origin":   url.Origin,
	}
	_, err := s.createStmt.ExecContext(ctx, params)
	return err
}

func (s urlService) Update(ctx context.Context, username, hash string, updates domain.UrlUpdates) error {
	params := map[string]any{
		"username":     username,
		"hash":         hash,
		"origin":       updates.Origin,
		"last_used_at": updates.LastUsedAt,
	}
	_, err := s.updateStmt.ExecContext(ctx, params)
	return err
}

func (s urlService) Delete(ctx context.Context, username, hash string) error {
	params := map[string]any{
		"username": username,
		"hash":     hash,
	}
	_, err := s.deleteStmt.ExecContext(ctx, params)
	return err
}
