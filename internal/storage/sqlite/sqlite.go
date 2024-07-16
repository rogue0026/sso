package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rogue0026/sso/internal/domain/models"
	"github.com/rogue0026/sso/internal/storage"
	"modernc.org/sqlite"
)

type Storage struct {
	pool *sql.DB
}

func New(dsn string) (*Storage, error) {
	const fn = "storage.sqlite.New"
	c, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	s := Storage{
		pool: c,
	}
	return &s, nil
}

func (s *Storage) Save(ctx context.Context, u models.User) (int64, error) {
	const fn = "storage.sqlite.Save"
	q := `INSERT INTO users (login, pass_hash, email) VALUES(?, ?, ?);`
	tx, err := s.pool.BeginTx(ctx, nil)
	if err != nil {
		return int64(0), fmt.Errorf("%s: %w", fn, err)
	}
	res, err := tx.ExecContext(ctx, q, u.Login, string(u.PassHash), u.Email)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == 19 {
			tx.Rollback()
			return int64(0), storage.ErrUserAlreadyExists
		}
		return int64(0), fmt.Errorf("%s: %w", fn, err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return int64(0), fmt.Errorf("%s: %w", fn, err)
	}
	tx.Commit()
	return lastId, nil
}

func (s *Storage) Fetch(ctx context.Context, login string, email string) (*models.User, error) {
	const fn = "storage.sqlite.Fetch"
	q := `SELECT login, pass_hash, email FROM users WHERE login = ? AND email = ?;`
	u := models.User{}
	if err := s.pool.QueryRowContext(ctx, q, login, email).Scan(u.Login, u.PassHash, u.Email); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &u, nil
}
