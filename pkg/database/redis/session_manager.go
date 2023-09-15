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

func (s sessionManagerService) Check(ctx context.Context, username string) (uuid.UUID, error) {
	value, err := db.Get(ctx, username).Result()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s sessionManagerService) Invalidate(ctx context.Context, username string) error {
	_, err := db.Pipelined(ctx, func(p redis.Pipeliner) error {
		id, err := p.GetDel(ctx, username).Result()
		if err != nil {
			return err
		}
		if err := p.Del(ctx, id).Err(); err != nil {
			return err
		}
		return nil
	})

	return err
}
