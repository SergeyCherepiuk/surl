package postgres

import (
	"context"

	"github.com/SergeyCherepiuk/surl/domain"
	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type verificationGetter struct {
	getStmt *sqlx.NamedStmt
}

func NewVerificationGetter() *verificationGetter {
	return &verificationGetter{
		getStmt: internal.MustPrepare(db, `SELECT * FROM verification_requests WHERE id = :id AND username = :username`),
	}
}

func (vg verificationGetter) Get(ctx context.Context, username string, id uuid.UUID) (domain.VerificationRequest, error) {
	params := map[string]any{
		"username": username,
		"id":       id,
	}
	var verificationRequest domain.VerificationRequest
	err := vg.getStmt.GetContext(ctx, &verificationRequest, params)
	return verificationRequest, err
}
