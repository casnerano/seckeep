package pgsql

import (
	"context"
	"errors"

	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/casnerano/seckeep/internal/server/repository"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository структура репозитория работы с записями пользователя.
type UserRepository struct {
	pgxpool *pgxpool.Pool
}

// NewUserRepository конструктор.
func NewUserRepository(pgxpool *pgxpool.Pool) repository.User {
	return &UserRepository{pgxpool}
}

// Add добавляет запись.
func (u UserRepository) Add(ctx context.Context, login, password, fullName string) (*model.User, error) {
	user := model.User{Login: login, Password: password}
	err := u.pgxpool.QueryRow(
		ctx,
		"insert into users(login, password, full_name) values($1, $2, $3) returning uuid, created_at",
		login,
		password,
		fullName,
	).Scan(
		&user.UUID,
		&user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = repository.ErrAlreadyExist
		}
		return nil, err
	}

	return &user, nil
}

// FindByLogin ищет запись по логину.
func (u UserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	user := model.User{Login: login}
	err := u.pgxpool.QueryRow(
		ctx,
		"select uuid, password, full_name, created_at from users where login = $1",
		login,
	).Scan(
		&user.UUID,
		&user.Password,
		&user.FullName,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

// FindByUUID ищет запись по UUID.
func (u UserRepository) FindByUUID(ctx context.Context, uuid string) (*model.User, error) {
	user := model.User{UUID: uuid}
	err := u.pgxpool.QueryRow(
		ctx,
		"select login, password, full_name, created_at from users where uuid = $1",
		uuid,
	).Scan(
		&user.Login,
		&user.Password,
		&user.FullName,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
