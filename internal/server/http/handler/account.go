package handler

//go:generate mockgen -destination=mock/account.go -source=account.go

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/casnerano/seckeep/internal/server/service/account"
	"github.com/casnerano/seckeep/pkg/log"
)

// AccountService интерфейс сервиса взаимодействия с аккаунтом.
type AccountService interface {
	SignIn(ctx context.Context, login, password string) (string, error)
	SignUp(ctx context.Context, login, password, fullName string) (string, error)
}

// Account структура обработчика взаимодействия с аккаунтом.
type Account struct {
	service AccountService
	logger  log.Loggable
}

// NewAccount конструктор.
func NewAccount(service AccountService, logger log.Loggable) *Account {
	return &Account{service: service, logger: logger}
}

// SignUp обработчик регистрации.
func (a *Account) SignUp(rd model.UserSignUpRequest, w http.ResponseWriter, r *http.Request) (any, int) {
	token, err := a.service.SignUp(r.Context(), rd.Login, rd.Password, rd.FullName)
	if err != nil {
		if errors.Is(err, account.ErrUserRegistered) {
			return nil, http.StatusConflict
		}

		a.logger.Error("Ошибка регистрации.", err)
		return nil, http.StatusInternalServerError
	}

	a.logger.Info(fmt.Sprintf("Успешная регистрация пользователя \"%s\"", rd.Login))
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return nil, http.StatusOK
}

// SignIn обработчик авторизации.
func (a *Account) SignIn(rd model.UserSignInRequest, w http.ResponseWriter, r *http.Request) (any, int) {
	token, err := a.service.SignIn(r.Context(), rd.Login, rd.Password)
	if err != nil {
		if errors.Is(err, account.ErrIncorrectCredentials) || errors.Is(err, account.ErrUserNotFound) {
			return nil, http.StatusUnauthorized
		}

		a.logger.Error("Ошибка авторизации.", err)
		return nil, http.StatusInternalServerError
	}

	a.logger.Info(fmt.Sprintf("Успешная авторизация пользователя. \"%s\"", rd.Login))
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return nil, http.StatusOK
}
