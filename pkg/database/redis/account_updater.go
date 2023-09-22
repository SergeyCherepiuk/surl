package redis

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/redis/go-redis/v9"
)

type accountUpdater struct {
	other domain.AccountUpdater
}

func NewAccountUpdater(other domain.AccountUpdater) *accountUpdater {
	return &accountUpdater{other: other}
}

func (au accountUpdater) UpdateUsername(ctx context.Context, username, newUsername string) error {
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

		return au.other.UpdateUsername(ctx, username, newUsername)
	})

	return err
}
