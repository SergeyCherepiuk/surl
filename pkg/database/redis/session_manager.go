package redis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

var SessionManagerService = sessionManagerService{}

type sessionManagerService struct{}

func (s sessionManagerService) Create(ctx context.Context, username string, ttl time.Duration) (uuid.UUID, error) {
	id := uuid.New()

	pipe := db.Pipeline()
	pipe.Set(ctx, username, id.String(), ttl)
	pipe.Set(ctx, id.String(), username, ttl)
	_, err := pipe.Exec(ctx)

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
	pipe := db.Pipeline()
	id, _ := pipe.GetDel(ctx, username).Result()
	pipe.Del(ctx, id)
	_, err := pipe.Exec(ctx)

	return err
}
