package domain

import (
	"context"
	"time"
)

type VerificationRequest struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type VerificationSender interface {
	Send(email string, verificationRequest VerificationRequest) error
}

type VerificationChecker interface {
	Check(ctx context.Context, username string) error
}

type VerificationGetter interface {
	Get(ctx context.Context, username, id string) (VerificationRequest, error)
}

type VerificationCreator interface {
	Create(ctx context.Context, verificationRequest VerificationRequest) error
}

type Verificator interface {
	Verify(ctx context.Context, username string) error
}

type VerificationDeleter interface {
	DeleteAll(ctx context.Context, username string) error
}
