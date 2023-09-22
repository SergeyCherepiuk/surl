package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type sessionCreator struct{}

func NewSessionCreator() *sessionCreator {
	return &sessionCreator{}
}

func (sc sessionCreator) Create(ctx context.Context, username string, ttl time.Duration) (uuid.UUID, error) {
	id := uuid.New()

	_, err := db.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.Set(ctx, username, id.String(), ttl).Err(); err != nil {
			return err
		}
		if err := p.Set(ctx, id.String(), username, ttl).Err(); err != nil {
			return err
		}
		return nil
	})

	return id, err
}
