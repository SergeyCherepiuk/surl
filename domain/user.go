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

type SessionManagerService interface {
	Create(ctx context.Context, username string, ttl time.Duration) (uuid.UUID, error)
	Check(ctx context.Context, username string) (uuid.UUID, error)
	Invalidate(ctx context.Context, username string) error
}

type AccountManagerService interface {
	Get(ctx context.Context, username string) (User, error)
	Create(ctx context.Context, user User) error
	Update(ctx context.Context, username string, updates map[string]any) error
	Delete(ctx context.Context, username string) error
}
