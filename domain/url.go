package domain

import (
	"context"
	"time"
)

type Url struct {
	Username   string    `json:"username" db:"username"`
	Hash       string    `json:"hash" db:"hash"`
	Origin     string    `json:"origin" db:"origin"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	LastUsedAt time.Time `json:"last_used_at" db:"last_used_at"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
}

type UrlUpdates struct {
	Origin     string
	Hash       string
	LastUsedAt time.Time
}

type OriginGetter interface {
	Get(ctx context.Context, username, hash string) (string, time.Duration, error)
}

type UrlGetter interface {
	Get(ctx context.Context, username, hash string) (Url, error)
	GetAll(ctx context.Context, username string) ([]Url, error)
	GetAllSorted(ctx context.Context, username, sortBy string, reversed bool) ([]Url, error)
}

type UrlCreator interface {
	Create(ctx context.Context, url Url) error
}

type UrlUpdater interface {
	Update(ctx context.Context, username, hash string, updates UrlUpdates) error
}

type UrlDeleter interface {
	Delete(ctx context.Context, username, hash string) error
}
