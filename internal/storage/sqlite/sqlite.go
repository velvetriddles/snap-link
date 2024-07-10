package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/velvetriddles/snap-link/internal/storage"
)

type Storage struct {
	db *sql.DB
}


func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON urls(alias);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO urls (alias, url) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"
	var url string
	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}
	return url, nil

}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"
	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
	}

	return nil
}
