package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type passwordResetCreator struct {
	createStmt *sqlx.NamedStmt
}

func NewPasswordResetCreator() *passwordResetCreator {
	return &passwordResetCreator{
		createStmt: internal.MustPrepare(db, `INSERT INTO password_reset_requests VALUES (:id, :username, :expires_at)`),
	}
}

func (prc passwordResetCreator) Create(ctx context.Context, passwordResetRequest domain.PasswordResetRequest) error {
	params := map[string]any{
		"id":         passwordResetRequest.ID,
		"username":   passwordResetRequest.Username,
		"expires_at": passwordResetRequest.ExpiresAt,
	}
	_, err := prc.createStmt.ExecContext(ctx, params)
	return err
}
