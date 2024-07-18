package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrBadSecretKey = errors.New("bad secret key")

type Claims struct {
	Login string
	jwt.RegisteredClaims
}

// NewJWT generates new jwt token based on user login and email
func NewJWT(login string) (string, error) {
	const fn = "lib.token.NewJWT"
	c := Claims{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &c, nil)

	// getting secret key for token signing
	k, err := key()
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	ss, err := token.SignedString([]byte(k))
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return ss, nil
}

// Validate validates signed token string. If token was corrupted, Validate returns error.
func Validate(token string) error {
	const fn = "lib.token.Validate"
	c := Claims{}

	_, err := jwt.ParseWithClaims(token, &c, func(t *jwt.Token) (interface{}, error) {
		k, err := key()
		if err != nil {
			return nil, err
		} else {
			return k, nil
		}
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func key() (string, error) {
	secret := os.Getenv("TOKEN_SIGNING_KEY")
	if len(secret) == 0 {
		return "", ErrBadSecretKey
	}

	return secret, nil
}
