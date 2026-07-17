package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenMaker struct {
	accessSecret  string
	refreshSecret string
}

func NewTokenMaker(accessSecret, refreshSecret string) (*TokenMaker, error) {
	if len(accessSecret) < 32 || len(refreshSecret) < 32 {
		return nil, errors.New("secret keys must be at least 32 characters long")
	}
	return &TokenMaker{accessSecret, refreshSecret}, nil
}

func (m *TokenMaker) CreateToken(userID uint64, duration time.Duration, isRefresh bool) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := m.accessSecret
	if isRefresh {
		secret = m.refreshSecret
	}

	return token.SignedString([]byte(secret))
}

func (m *TokenMaker) VerifyToken(tokenString string, isRefresh bool) (*TokenClaims, error) {
	secret := m.accessSecret
	if isRefresh {
		secret = m.refreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}