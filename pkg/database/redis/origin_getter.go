package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/redis/internal"
)

type originGetter struct {
	other domain.OriginGetter
}

func NewOriginGetter(other domain.OriginGetter) *originGetter {
	return &originGetter{other: other}
}

func (og originGetter) Get(ctx context.Context, username, hash string) (string, time.Duration, error) {
	key := fmt.Sprintf("%s/%s", username, hash)

	if origin, err := cacheDb.Get(ctx, key).Result(); err == nil {
		ttl, _ := cacheDb.TTL(ctx, key).Result()
		return origin, ttl, nil
	}

	origin, expiresIn, err := og.other.Get(context.Background(), username, hash)
	ttl := internal.Min(5*time.Minute, expiresIn)
	if err == nil {
		cacheDb.Set(ctx, key, origin, ttl)
	}

	return origin, ttl, err
}
