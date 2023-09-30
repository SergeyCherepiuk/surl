package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/SergeyCherepiuk/surl/domain"
)

type urlUpdater struct {
	other domain.UrlUpdater
}

func NewUrlUpdater(other domain.UrlUpdater) *urlUpdater {
	return &urlUpdater{other: other}
}

func (uu urlUpdater) Update(ctx context.Context, username, hash string, updates domain.UrlUpdates) error {
	// Will error (won't be found) if hash (and/or origin) is changed
	if err := cacheDb.Get(ctx, fmt.Sprintf("%s/%s", username, updates.Hash)).Err(); err != nil {
		// Get TTL of the old cache
		ttl, _ := cacheDb.TTL(ctx, fmt.Sprintf("%s/%s", username, hash)).Result()

		// Delete old cache
		if err := cacheDb.Del(ctx, fmt.Sprintf("%s/%s", username, hash)).Err(); err != nil {
			return err
		}

		// Cache new link if the old cache was still valid
		if ttl != -2*time.Nanosecond {
			cacheDb.Set(ctx, fmt.Sprintf("%s/%s", username, updates.Hash), updates.Origin, ttl) // NOTE: Error is ignored
		}
	}

	return uu.other.Update(context.Background(), username, hash, updates)
}
