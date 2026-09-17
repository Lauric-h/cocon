package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS meals (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	date       TEXT NOT NULL,
	meal_slot  TEXT NOT NULL CHECK(meal_slot IN ('petit-déjeuner', 'midi', 'dîner')),
	dish       TEXT NOT NULL,
	notes      TEXT DEFAULT '',
	UNIQUE(date, meal_slot)
);

CREATE TABLE IF NOT EXISTS groceries (
	id         INTEGER PRIMARY KEY AUTOINCREMENT,
	name       TEXT NOT NULL,
	quantity   INTEGER, -- NULL si non précisé
	checked    INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL
);
`

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("migrate schema: %w", err)
	}

	return db, nil
}
