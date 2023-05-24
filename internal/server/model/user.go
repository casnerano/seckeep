package model

import "time"

// User структура пользователя (представляет модель БД).
type User struct {
	UUID      string
	Login     string
	Password  string
	FullName  string
	CreatedAt time.Time
}

// UserSignUpRequest структура запроса на регистрацию.
type UserSignUpRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}

// UserSignInRequest структура запроса на авторизацию.
type UserSignInRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
