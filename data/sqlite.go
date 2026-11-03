package data

import (
	"database/sql"
	"fmt"

	"github.com/okunix/ambicertd/errutil"
	_ "modernc.org/sqlite"
)

func SQLite(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, errutil.WithAttrs(fmt.Errorf("failed to open db connection: %w", err), "dsn", dsn)
	}
	if err := db.Ping(); err != nil {
		return nil, errutil.WithAttrs(fmt.Errorf("failed to ping db: %w", err), "dsn", dsn)
	}
	return db, nil
}
