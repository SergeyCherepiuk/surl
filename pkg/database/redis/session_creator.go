package redis

import (
	"context"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type sessionCreator struct {
	accountCreator domain.AccountCreator
}

func NewSessionCreator(accountCreator domain.AccountCreator) *sessionCreator {
	return &sessionCreator{accountCreator: accountCreator}
}

func (sc sessionCreator) Create(ctx context.Context, user domain.User, ttl time.Duration) (uuid.UUID, error) {
	id := uuid.New()

	_, err := db.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.Set(ctx, user.Username, id.String(), ttl).Err(); err != nil {
			return err
		}
		if err := p.Set(ctx, id.String(), user.Username, ttl).Err(); err != nil {
			return err
		}

		return sc.accountCreator.Create(ctx, user)
	})

	return id, err
}
