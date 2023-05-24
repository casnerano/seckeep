// Package repository содержит наборы методов для взаимодействия с сущностями в хранилище.
package repository

//go:generate mockgen -destination=mock/repository.go -source=repository.go

import "errors"

import (
	"context"
	"time"

	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/casnerano/seckeep/internal/shared"
)

// Общие ошибки при работы БД.
var (
	// ErrAlreadyExist Запись существует.
	ErrAlreadyExist = errors.New("already exists")

	// ErrNotFound Запись не найдена.
	ErrNotFound = errors.New("not found")
)

// User интерфейс работы с записями пользователей.
type User interface {
	// Add добавляет запись.
	Add(ctx context.Context, login, password, fullName string) (*model.User, error)

	// FindByLogin ищет запись по логину.
	FindByLogin(ctx context.Context, login string) (*model.User, error)

	// FindByUUID ищет запись по UUID.
	FindByUUID(ctx context.Context, uuid string) (*model.User, error)
}

// Data интерфейс работы с записями секретных данных.
type Data interface {
	// Add добавляет запись.
	Add(ctx context.Context, data shared.Data) (*shared.Data, error)

	// FindByUUID ищет запись по UUID.
	FindByUUID(ctx context.Context, userUUID string, uuid string) (*shared.Data, error)

	// FindByUserUUID ищет запись по UUID пользователя.
	FindByUserUUID(ctx context.Context, userUUID string) ([]*shared.Data, error)

	// Update обновляет запись.
	Update(ctx context.Context, userUUID string, uuid string, value []byte, version time.Time) (*shared.Data, error)

	// Delete удаляет запись.
	Delete(ctx context.Context, userUUID string, uuid string) error
}
