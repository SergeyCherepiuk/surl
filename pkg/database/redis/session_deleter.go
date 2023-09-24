package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type sessionDeleter struct{}

func NewSessionDeleter() *sessionDeleter {
	return &sessionDeleter{}
}

func (ad sessionDeleter) Delete(ctx context.Context, username string) error {
	id, err := db.Get(ctx, username).Result()
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
