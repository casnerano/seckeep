// Package jwtoken содердит методы работы с JWT токеном.
package jwtoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Основные ошибки при работе с JWT.
var (
	// ErrInvalidToken некорректный токен.
	ErrInvalidToken = errors.New("invalid token")
)

// Payload структура полезных данных.
type Payload struct {
	UUID     string
	FullName string
}

// Claims структура содержимого токена.
type Claims struct {
	Payload Payload
	jwt.RegisteredClaims
}

// JWToken структура токена.
type JWToken struct{}

// New конструктор.
func New() *JWToken {
	return &JWToken{}
}

// Create метод создания токена.
func (j JWToken) Create(payload Payload, ttl time.Duration, secret []byte) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			payload,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	)
	ss, err := token.SignedString(secret)
	if err != nil {
		return "", nil
	}
	return ss, nil
}

// Verify метод верификации токена.
func (j JWToken) Verify(token string, secret []byte) (*Payload, error) {
	claims := Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, ErrInvalidToken
	}

	return &claims.Payload, nil
}
