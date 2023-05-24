package account

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/go-resty/resty/v2"
)

// Основные ошибки авторизации.
var (
	// ErrIncorrectCredentials некорректные учетные данные.
	ErrIncorrectCredentials = errors.New("incorrect credentials")

	// ErrUserRegistered пользователь зарегистрирован.
	ErrUserRegistered = errors.New("user is registered")
)

// Account структура для авторизации и регистрации пользователя на сервере.
type Account struct {
	client   *resty.Client
	tokenJar *TokenJar
}

// New конструктор.
func New(client *resty.Client, tokenJar *TokenJar) *Account {
	return &Account{
		client:   client,
		tokenJar: tokenJar,
	}
}

// SignUp метод регистрации на сервере.
func (a Account) SignUp(login, password, fullName string) error {
	body := model.UserSignUpRequest{
		Login:    login,
		Password: password,
		FullName: fullName,
	}
	response, err := a.client.R().SetBody(body).Post("/user/register")
	if err != nil {
		return err
	}

	switch response.StatusCode() {
	case http.StatusBadRequest:
		return fmt.Errorf("incorrect values: %w", errors.New(string(response.Body())))
	case http.StatusConflict:
		return ErrUserRegistered
	case http.StatusOK:
		return a.flushHeaderToken(response.Header())
	}

	return fmt.Errorf("internal server error: %w", errors.New(response.Status()))
}

// SignIn метод авторизации на сервере.
func (a Account) SignIn(login, password string) error {
	body := model.UserSignInRequest{Login: login, Password: password}
	response, err := a.client.R().SetBody(body).Post("/user/login")
	if err != nil {
		return err
	}

	switch response.StatusCode() {
	case http.StatusBadRequest:
		return fmt.Errorf("incorrect values: %w", errors.New(string(response.Body())))
	case http.StatusUnauthorized:
		return ErrIncorrectCredentials
	case http.StatusOK:
		return a.flushHeaderToken(response.Header())
	}

	return fmt.Errorf("internal server error: %w", errors.New(response.Status()))
}

// flushHeaderToken метод сбрасывает (сохраняет) токен из заголовков.
func (a Account) flushHeaderToken(header http.Header) error {
	parts := strings.Split(header.Get("Authorization"), " ")
	if len(parts) != 2 {
		return errors.New("authorization token not found")
	}
	if err := a.tokenJar.SetToken(parts[1]); err != nil {
		return fmt.Errorf("token jar error: %w", err)
	}
	return nil
}
