package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rogue0026/sso/internal/domain/models"
	"github.com/rogue0026/sso/internal/lib/token"
	"github.com/rogue0026/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidUserCredentials = errors.New("invalid user credentials")

type UserSaver interface {
	Save(ctx context.Context, u models.User) (int64, error)
}

type UserFetcher interface {
	Fetch(ctx context.Context, login string, email string) (models.User, error)
}

type Service struct {
	Logger  *slog.Logger
	Saver   UserSaver
	Fetcher UserFetcher
}

func New(l *slog.Logger, s UserSaver, f UserFetcher) *Service {
	l = l.With("layer", "service")
	svc := Service{
		Logger:  l,
		Saver:   s,
		Fetcher: f,
	}
	return &svc
}

func (s *Service) RegisterNewUser(ctx context.Context, login string, password string, email string) (int64, error) {
	const fn = "services.auth.RegisterNewUser"
	l := s.Logger.With("func", fn)
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		l.Error(err.Error())
		return int64(0), fmt.Errorf("%s:%w", fn, err)
	}
	userId, err := s.Saver.Save(ctx,
		models.User{
			Login:    login,
			PassHash: passHash,
			Email:    email,
		})
	if err != nil {
		l.Error(err.Error())
		return int64(0), fmt.Errorf("%s:%w", fn, err)
	}
	return userId, nil
}

func (s *Service) LoginUser(ctx context.Context, login string, password string, email string) (string, error) {
	const fn = "services.auth.LoginUser"
	l := s.Logger.With("func", fn)
	// searching user in database
	user, err := s.Fetcher.Fetch(ctx, login, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", storage.ErrUserNotFound
		} else {
			l.Error(err.Error())
			return "", fmt.Errorf("%s: %w", fn, err)
		}
	}

	// check user credentials
	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		l.Error(err.Error())
		return "", ErrInvalidUserCredentials
	}

	// generate jwt token and send it back to client
	tokenStr, err := token.NewJWT(login, email)
	if err != nil {
		l.Error(err.Error())
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return tokenStr, nil
}
