package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rogue0026/sso/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	Save(u models.User) (int64, error)
}

type UserFetcher interface {
	Fetch(login string, email string) (models.User, error)
}

type Service struct {
	Logger  *slog.Logger
	Saver   UserSaver
	Fetcher UserFetcher
}

func (s *Service) RegisterNewUser(ctx context.Context, login string, password string, email string) (int64, error) {
	const fn = "services.auth.RegisterNewUser"
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return int64(0), fmt.Errorf("%s:%w", fn, err)
	}
	userId, err := s.Saver.Save(
		models.User{
			Login:    login,
			PassHash: passHash,
			Email:    email,
		})
	if err != nil {
		return int64(0), fmt.Errorf("%s:%w", fn, err)
	}
	return userId, nil
}

func (s *Service) LoginUser(ctx context.Context, login string, password string) (string, error) {

	return "", nil
}
