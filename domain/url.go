package domain

import (
	"context"
	"hash"
)

type Url struct {
	Username string      `json:"username" db:"username"`
	Hash     hash.Hash32 `json:"hash" db:"hash"`
	Origin   string      `json:"origin" db:"origin"`
}

type UrlService interface {
	GetAll(ctx context.Context, username string) ([]Url, error)
	GetOrigin(ctx context.Context, username, hash string) (string, error)
	Create(ctx context.Context, username, origin string) error
	Update(ctx context.Context, username, hash, newOrigin string) error
	Delete(ctx context.Context, username, hash string) error
}
