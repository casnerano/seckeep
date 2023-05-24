// Package data содержит методы для работы с секретными данными пользователя.
package data

import (
	"context"
	"errors"
	"time"

	"github.com/casnerano/seckeep/internal/server/repository"
	"github.com/casnerano/seckeep/internal/shared"
)

// Основные ошибки при работе с секретными данными пользователя..
var (
	// ErrNotFound запись не найдена.
	ErrNotFound = errors.New("not found")
)

// Data структура для работы с секретными данными пользователя.
type Data struct {
	repo repository.Data
}

// New конструктор.
func New(repo repository.Data) *Data {
	return &Data{
		repo: repo,
	}
}

// Create метод для создания.
func (d Data) Create(ctx context.Context, data shared.Data) (*shared.Data, error) {
	return d.repo.Add(ctx, data)
}

// FindByUUID метод поиска по UUID.
func (d Data) FindByUUID(ctx context.Context, userUUID, uuid string) (*shared.Data, error) {
	data, err := d.repo.FindByUUID(ctx, userUUID, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

// FindByUserUUID метод по UUID пользователя.
func (d Data) FindByUserUUID(ctx context.Context, userUUID string) ([]*shared.Data, error) {
	return d.repo.FindByUserUUID(ctx, userUUID)
}

// Update метод обновления.
func (d Data) Update(ctx context.Context, userUUID, uuid string, value []byte, version time.Time) (*shared.Data, error) {
	data, err := d.repo.Update(ctx, userUUID, uuid, value, version)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return data, nil
}

// Delete метод удаления.
func (d Data) Delete(ctx context.Context, userUUID, uuid string) error {
	err := d.repo.Delete(ctx, userUUID, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
