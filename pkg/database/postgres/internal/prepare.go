package internal

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func MustPrepare(db *sqlx.DB, query string) *sqlx.NamedStmt {
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}
