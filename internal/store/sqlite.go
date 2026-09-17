package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS tasks (
	id                      INTEGER PRIMARY KEY AUTOINCREMENT,
	title                   TEXT NOT NULL,
	description             TEXT DEFAULT '',
	category                TEXT NOT NULL,
	due_date                TEXT,    -- NULL si backlog ou récurrente
	recurrence_frequency    TEXT,    -- daily | weekly | monthly | yearly
	recurrence_interval     INTEGER,
	recurrence_days_of_week TEXT,    -- ex: "1,3" (lundi=1, mercredi=3)
	assigned_to             TEXT DEFAULT '',
	source                  TEXT NOT NULL DEFAULT 'manual'
);

CREATE TABLE IF NOT EXISTS task_occurrences (
	id           INTEGER PRIMARY KEY AUTOINCREMENT,
	task_id      INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	date         TEXT NOT NULL,
	assigned_to  TEXT DEFAULT '',
	done         INTEGER NOT NULL DEFAULT 0,
	completed_at TEXT,
	UNIQUE(task_id, date)
);

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
