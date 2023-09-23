package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB //db connection
}

func NewSqlite(config map[string]string) (*Storage, error) {
	db, err := sql.Open("sqlite3", config["sqlite_path"])
	if err != nil {
		return nil, fmt.Errorf("cant open connection, %s", err)
	}
	if err = createDb(db); err != nil {
		return nil, fmt.Errorf("cant create database, %s", err)
	}

	return &Storage{db}, nil
}

func createDb(db *sql.DB) error {
	const op = "storage.sqlite.createDb"

	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
		    id INTEGER PRIMARY KEY,
		    origin TEXT NOT NULL,
		    alias TEXT NOT NULL
		)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = statement.Exec()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Save(originUrl string, alias string) (int64, error) {
	const op = "storage.sqlite.save"

	stm, err := s.db.Prepare("INSERT INTO url (origin, alias) VALUES (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stm.Exec(originUrl, alias)

	if err != nil {
		if sqlErr, ok := err.(sqlite3.Error); ok && sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, UrlExist)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetByAlias(alias string) (*UrlModel, error) {
	const op = "storage.sqlite.getByOrigin"

	stm, err := s.db.Prepare("SELECT id, origin, alias FROM url where alias = ?")

	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement %w", op, err)
	}

	var url UrlModel
	err = stm.QueryRow(alias).Scan(&url.id, &url.origin, &url.alias)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		}

		return nil, fmt.Errorf("%s: execute statement %w", op, err)
	}

	return &url, nil
}

func (s *Storage) DeleteByAlias(alias string) error {
	const op = "storage.sqlite.deleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")

	if err != nil {
		return fmt.Errorf("%s: prepare statement %w", op, err)
	}

	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s: execute statement %w", op, err)
	}

	return nil
}
