package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type accountUpdater struct {
	updateUsernameStmt *sqlx.NamedStmt
	updatePasswordStmt *sqlx.NamedStmt
}

func NewAccountUpdater() *accountUpdater {
	return &accountUpdater{
		updateUsernameStmt: internal.MustPrepare(db, `UPDATE users SET username = :new_username WHERE username = :username`),
		updatePasswordStmt: internal.MustPrepare(db, `UPDATE users SET password = :new_password WHERE username = :username`),
	}
}

func (au accountUpdater) UpdateUsername(ctx context.Context, username, newUsername string) error {
	params := map[string]any{
		"username":     username,
		"new_username": newUsername,
	}
	_, err := au.updateUsernameStmt.ExecContext(ctx, params)
	return err
}

func (au accountUpdater) UpdatePassword(ctx context.Context, username, newPassword string) error {
	params := map[string]any{
		"username": username,
		"new_password": newPassword,
	}
	_, err := au.updatePasswordStmt.ExecContext(ctx, params)
	return err
}