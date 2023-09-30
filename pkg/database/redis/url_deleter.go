package redis

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/surl/domain"
)

type urlDeleter struct {
	other domain.UrlDeleter
}

func NewUrlDeleter(other domain.UrlDeleter) *urlDeleter {
	return &urlDeleter{other: other}
}

func (ud urlDeleter) Delete(ctx context.Context, username, hash string) error {
	if err := cacheDb.Del(ctx, fmt.Sprintf("%s/%s", username, hash)).Err(); err != nil {
		return err
	}

	return ud.other.Delete(ctx, username, hash)
}
