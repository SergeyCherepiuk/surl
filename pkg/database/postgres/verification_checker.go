package postgres

import (
	"context"
	"fmt"

	"github.com/SergeyCherepiuk/surl/pkg/database/postgres/internal"
	"github.com/jmoiron/sqlx"
)

type verificationChecker struct {
	checkStmt *sqlx.NamedStmt
}

func NewVerificationChecker() *verificationChecker {
	return &verificationChecker{
		checkStmt: internal.MustPrepare(db, `SELECT is_verified FROM users WHERE username = :username`),
	}
}

func (vc verificationChecker) Check(ctx context.Context, username string) error {
	params := map[string]any{"username": username}
	var isVerified bool
	err := vc.checkStmt.GetContext(ctx, &isVerified, params)
	if err != nil || !isVerified {
		return fmt.Errorf("account is not verified")
	}
	return nil
}
