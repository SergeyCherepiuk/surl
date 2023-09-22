package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type SessionChecker interface {
	Check(ctx context.Context, id uuid.UUID) (string, error)
}

type AccountGetter interface {
	Get(ctx context.Context, username string) (User, error)
}

type SessionCreator interface {
	Create(ctx context.Context, username string, ttl time.Duration) (uuid.UUID, error)
}

type AccountCreator interface {
	Create(ctx context.Context, user User) error
}

type SessionUpdater interface {
	UpdateUsername(ctx context.Context, username, newUsername string) error
}

type AccountUpdater interface {
	UpdateUsername(ctx context.Context, username, newUsername string) error
	UpdatePassword(ctx context.Context, username, newPassword string) error
}

type AccountDeleter interface {
	Delete(ctx context.Context, username string) error
}
