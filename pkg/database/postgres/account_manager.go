package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type accountManagerService struct {
	getStmt    *sqlx.NamedStmt
	createStmt *sqlx.NamedStmt
	updateStmt *sqlx.NamedStmt
	deleteStmt *sqlx.NamedStmt
}

func NewAccountManagerService() *accountManagerService {
	return &accountManagerService{
		getStmt:    internal.MustPrepare(db, `SELECT * FROM users WHERE username = :username`),
		createStmt: internal.MustPrepare(db, `INSERT INTO users VALUES (:username, :password)`),
		updateStmt: internal.MustPrepare(db, `UPDATE users SET username = :new_username, password = :new_password WHERE username = :username`),
		deleteStmt: internal.MustPrepare(db, `DELETE FROM users WHERE username = :username`),
	}
}

func (s accountManagerService) Get(ctx context.Context, username string) (domain.User, error) {
	params := map[string]any{"username": username}
	user := domain.User{}
	err := s.getStmt.GetContext(ctx, &user, params)
	return user, err
}

func (s accountManagerService) Create(ctx context.Context, user domain.User) error {
	params := map[string]any{
		"username": user.Username,
		"password": user.Password,
	}
	_, err := s.createStmt.ExecContext(ctx, params)
	return err
}

func (s accountManagerService) Update(ctx context.Context, username string, updates map[string]any) error {
	user, err := s.Get(ctx, username)
	if err != nil {
		return err
	}

	newUsername := user.Username
	if value, ok := updates["username"].(string); ok && newUsername != value {
		newUsername = value
	}
	newPassword := user.Password
	if value, ok := updates["password"].(string); ok && newPassword != value {
		newPassword = value
	}

	params := map[string]any{
		"new_username": newUsername,
		"new_password": newPassword,
		"username":     username,
	}
	_, err = s.updateStmt.ExecContext(ctx, params)
	return err
}

func (s accountManagerService) Delete(ctx context.Context, username string) error {
	params := map[string]any{"username": username}
	_, err := s.deleteStmt.ExecContext(ctx, params)
	return err
}
