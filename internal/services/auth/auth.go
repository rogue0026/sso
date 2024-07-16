package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rogue0026/sso/internal/domain/models"
	"github.com/rogue0026/sso/internal/lib/token"
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

func (s *Service) RegisterNewUser(ctx context.Context, login string, password string, email string) (int64, error) {
	const fn = "services.auth.RegisterNewUser"
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return int64(0), fmt.Errorf("%s:%w", fn, err)
	}
	userId, err := s.Saver.Save(ctx,
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

func (s *Service) LoginUser(ctx context.Context, login string, password string, email string) (string, error) {
	const fn = "services.auth.LoginUser"
	// 1. проверить, есть ли в базе данных пользователь с указанным логином и если есть, то проверить правильность введенного пароля
	// 2. Если пароль введен неправильно, то вернуть ошибку
	// 3. Если пароль введен правильно, то сгенерировать jwt-токен и вернуть его пользователю
	user, err := s.Fetcher.Fetch(ctx, login, email)
	if err != nil {
		// todo
		// user not found error
		// internal error in persistent layer
		return "", fmt.Errorf("%s: %w", fn, err)
	}
	// check user credentials
	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		// send empty string and invalid user credentials error
		return "", ErrInvalidUserCredentials
	}

	// generate jwt token and send it back to client
	tokenStr, err := token.NewJWT(login, email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return tokenStr, nil
}
