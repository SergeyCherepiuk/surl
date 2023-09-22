package redis

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/redis/go-redis/v9"
)

type accountDeleter struct {
	other domain.AccountDeleter
}

func NewAccountDeleter(other domain.AccountDeleter) *accountDeleter {
	return &accountDeleter{other: other}
}

func (ad accountDeleter) Delete(ctx context.Context, username string) error {
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

		return ad.other.Delete(ctx, username)
	})

	return err
}
