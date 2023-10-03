package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type passwordResetGetter struct {
	getStmt *sqlx.NamedStmt
}

func NewPasswordResetGetter() *passwordResetGetter {
	return &passwordResetGetter{
		getStmt: internal.MustPrepare(db, `SELECT * FROM password_reset_requests WHERE username = :username AND id = :id`),
	}
}

func (prg passwordResetGetter) Get(ctx context.Context, username, id string) (domain.PasswordResetRequest, error) {
	params := map[string]any{
		"username": username,
		"id":       id,
	}
	var passwordResetRequest domain.PasswordResetRequest
	err := prg.getStmt.GetContext(ctx, &passwordResetRequest, params)
	return passwordResetRequest, err
}
