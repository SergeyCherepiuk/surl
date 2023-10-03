package domain

import (
	"context"
	"time"
)

type PasswordResetRequest struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type PasswordResetSender interface {
	Send(ctx context.Context, email string, passwordResetRequest PasswordResetRequest) error
}

type PasswordResetGetter interface {
	Get(ctx context.Context, username, id string) (PasswordResetRequest, error)
}

type PasswordResetCreator interface {
	Create(ctx context.Context, passwordResetRequest PasswordResetRequest) error
}

type PasswordResetDeleter interface {
	DeleteAll(ctx context.Context, username string) error
}
