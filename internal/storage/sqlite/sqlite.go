package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rogue0026/sso/internal/domain/models"
	"github.com/rogue0026/sso/internal/storage"
	"modernc.org/sqlite"
)

type Storage struct {
	logger *slog.Logger
	pool   *sql.DB
}

func New(l *slog.Logger, dsn string) (*Storage, error) {
	const fn = "storage.sqlite.New"
	l = l.With("layer", "database")
	c, err := sql.Open("sqlite", dsn)
	if err != nil {
		l.Error(err.Error(), "func", fn)
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	s := Storage{
		logger: l,
		pool:   c,
	}
	return &s, nil
}

func (s *Storage) Save(ctx context.Context, u models.User) (int64, error) {
	const fn = "storage.sqlite.Save"
	l := s.logger.With("func", fn)
	q := `INSERT INTO users (login, pass_hash, email) VALUES(?, ?, ?);`
	tx, err := s.pool.BeginTx(ctx, nil)
	if err != nil {
		l.Error(err.Error())
		return int64(0), fmt.Errorf("%s: %w", fn, err)
	}
	res, err := tx.ExecContext(ctx, q, u.Login, string(u.PassHash), u.Email)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == 19 {
			l.Error(err.Error())
			tx.Rollback()
			return int64(0), storage.ErrUserAlreadyExists
		}
		l.Error(err.Error())
		return int64(0), fmt.Errorf("%s: %w", fn, err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		l.Error(err.Error())
		return int64(0), fmt.Errorf("%s: %w", fn, err)
	}
	tx.Commit()
	return lastId, nil
}

func (s *Storage) Fetch(ctx context.Context, login string, email string) (*models.User, error) {
	const fn = "storage.sqlite.Fetch"
	l := s.logger.With("func", fn)
	q := `SELECT login, pass_hash, email FROM users WHERE login = ? AND email = ?;`
	u := models.User{}
	if err := s.pool.QueryRowContext(ctx, q, login, email).Scan(u.Login, u.PassHash, u.Email); err != nil {
		l.Error(err.Error())
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &u, nil
}
