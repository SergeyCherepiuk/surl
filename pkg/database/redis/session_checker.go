package redis

import (
	"context"

	"github.com/google/uuid"
)

type sessionChecker struct{}

func NewSessionChecker() *sessionChecker {
	return &sessionChecker{}
}

func (sc sessionChecker) Check(ctx context.Context, id uuid.UUID) (string, error) {
	return db.Get(ctx, id.String()).Result()
}
