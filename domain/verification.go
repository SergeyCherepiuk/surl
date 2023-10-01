package domain

import (
	"context"

	"github.com/google/uuid"
)

type VerificationRequest struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
}

type VerificationChecker interface {
	Check(ctx context.Context, username string) error
}

type VerificationGetter interface {
	Get(ctx context.Context, username string, id uuid.UUID) (VerificationRequest, error)
}

type Verificator interface {
	Verify(ctx context.Context, username string) error
}

type VerificationDeleter interface {
	DeleteAll(ctx context.Context, username string) error
}
