package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type verificationCreator struct {
	createStmt *sqlx.NamedStmt
}

func NewVerificationCreator() *verificationCreator {
	return &verificationCreator{
		createStmt: internal.MustPrepare(db, `INSERT INTO verification_requests VALUES (:id, :username, :expires_at)`),
	}
}

func (vc verificationCreator) Create(ctx context.Context, verificationRequest domain.VerificationRequest) error {
	params := map[string]any{
		"username":   verificationRequest.Username,
		"id":         verificationRequest.ID,
		"expires_at": verificationRequest.ExpiresAt,
	}
	_, err := vc.createStmt.ExecContext(ctx, params)
	return err
}
