package domain

import "hash"

type Url struct {
	Username string      `json:"username" db:"username"`
	Hash     hash.Hash32 `json:"hash" db:"hash"`
	Origin   string      `json:"origin" db:"origin"`
}

type UrlService interface {
	GetAll(username string) ([]Url, error)
	GetOrigin(username, hash string) (string, error)
	Create(username, origin string) (Url, error)
	Update(username, hash, newOrigin string) (Url, error)
	Delete(username, hash string) (Url, error)
}
