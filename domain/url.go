package domain

import "hash"

type Url struct {
	Username string      `json:"username" db:"username"`
	Hash     hash.Hash32 `json:"hash" db:"hash"`
	Origin   string      `json:"origin" db:"origin"`
}
