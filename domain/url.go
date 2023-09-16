package domain

import (
	"context"
)

type Url struct {
	Username string `json:"username" db:"username"`
	Hash     string `json:"hash" db:"hash"`
	Origin   string `json:"origin" db:"origin"`
}

type UrlService interface {
	GetOrigin(ctx context.Context, username, hash string) (string, error)
	GetAll(ctx context.Context, username string) ([]Url, error)
	Create(ctx context.Context, url Url) error
	Update(ctx context.Context, username, hash, newOrigin string) error
	Delete(ctx context.Context, username, hash string) error
}
