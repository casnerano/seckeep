// Package account содержит методы работы с аккантом пользователя.
package account

//go:generate mockgen -destination=mock/account.go -source=account.go

import (
	"context"
	"errors"
	"time"

	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/casnerano/seckeep/internal/server/repository"
	"github.com/casnerano/seckeep/pkg/jwtoken"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtTTL = 15 * time.Minute
)

// Основные ошибки при работе с аккантом.
var (
	// ErrIncorrectCredentials неверные очетные даные.
	ErrIncorrectCredentials = errors.New("incorrect credentials")

	// ErrUserRegistered пользователь зарегистрирован.
	ErrUserRegistered = errors.New("user is registered")

	// ErrUserNotFound пользователь не найден.
	ErrUserNotFound = errors.New("user not found")
)

// JWT интерфейс работы с JWT токеном.
type JWT interface {
	Create(payload jwtoken.Payload, ttl time.Duration, secret []byte) (string, error)
}

// Account структура для работы с аккантом пользователя.
type Account struct {
	repo   repository.User
	jwt    JWT
	secret string
}

// New конструктор.
func New(repo repository.User, jwt JWT, secret string) *Account {
	return &Account{
		repo:   repo,
		jwt:    jwt,
		secret: secret,
	}
}

// SignIn метод авторизации.
func (a Account) SignIn(ctx context.Context, login, password string) (string, error) {
	user, err := a.repo.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrIncorrectCredentials
	}

	return a.createTokenForUser(user)
}

// SignUp метод регистрации.
func (a Account) SignUp(ctx context.Context, login, password, fullName string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := a.repo.Add(ctx, login, string(hashedPassword), fullName)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExist) {
			return "", ErrUserRegistered
		}
		return "", err
	}

	return a.createTokenForUser(user)
}

// createTokenForUser метод генерирует токен для заданного пользователя.
func (a Account) createTokenForUser(user *model.User) (string, error) {
	return a.jwt.Create(jwtoken.Payload{UUID: user.UUID, FullName: user.FullName}, jwtTTL, []byte(a.secret))
}
