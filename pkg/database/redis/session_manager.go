package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type sessionManagerService struct{}

func NewSessionManagerService() *sessionManagerService {
	return &sessionManagerService{}
}

func (s sessionManagerService) Create(ctx context.Context, username string, ttl time.Duration) (uuid.UUID, error) {
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

func (s sessionManagerService) Check(ctx context.Context, id uuid.UUID) (string, error) {
	return db.Get(ctx, id.String()).Result()
}

func (s sessionManagerService) Invalidate(ctx context.Context, username string) error {
	id, err := db.Get(context.Background(), username).Result()
	if err != nil {
		return err
	}

	_, err = db.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.Del(ctx, username).Err(); err != nil {
			return err
		}
		if err := p.Del(ctx, id).Err(); err != nil {
			return err
		}
		return nil
	})

	return err
}
