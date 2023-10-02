package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type verificationCreator struct {
	createStmt *sqlx.NamedStmt
}

func NewVerificationCreator() *verificationCreator {
	return &verificationCreator{
		createStmt: internal.MustPrepare(db, `INSERT INTO verification_requests VALUES (:id, :username)`),
	}
}

func (vc verificationCreator) Create(ctx context.Context, username, id string) error {
	params := map[string]any{
		"username": username,
		"id":       id,
	}
	_, err := vc.createStmt.ExecContext(ctx, params)
	return err
}
