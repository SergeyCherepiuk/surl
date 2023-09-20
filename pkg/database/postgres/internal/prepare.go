package internal

import (
	"fmt"
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

func MustPrepareMap(db *sqlx.DB, keys []string, query string) map[string]*sqlx.NamedStmt {
	stmts := make(map[string]*sqlx.NamedStmt, len(keys))
	for _, key := range keys {
		stmts[key] = MustPrepare(db, fmt.Sprintf(query, key))
	}
	return stmts
}
