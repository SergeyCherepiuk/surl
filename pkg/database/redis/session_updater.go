package redis

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/redis/go-redis/v9"
)

type sessionUpdater struct {
	accountUpdater domain.AccountUpdater
}

func NewAccountUpdater(accountUpdater domain.AccountUpdater) *sessionUpdater {
	return &sessionUpdater{accountUpdater: accountUpdater}
}

func (au sessionUpdater) UpdateUsername(ctx context.Context, username, newUsername string) error {
	id, err := db.Get(ctx, username).Result()
	if err != nil {
		return err
	}

	ttl, err := db.TTL(ctx, username).Result()
	if err != nil {
		return err
	}

	_, err = db.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.Set(ctx, id, newUsername, ttl).Err(); err != nil {
			return err
		}
		if err := p.Set(ctx, newUsername, id, ttl).Err(); err != nil {
			return err
		}
		if err := p.Del(ctx, username).Err(); err != nil {
			return err
		}

		return au.accountUpdater.UpdateUsername(ctx, username, newUsername)
	})

	return err
}
