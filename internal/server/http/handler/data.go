package handler

//go:generate mockgen -destination=mock/data.go -source=data.go

import (
	"context"
	"errors"
	"net/http"
	"time"

	smodel "github.com/casnerano/seckeep/internal/pkg/model"
	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/go-chi/chi/v5"

	"github.com/casnerano/seckeep/internal/server/http/middleware"
	"github.com/casnerano/seckeep/internal/server/service/data"
	"github.com/casnerano/seckeep/pkg/log"
)

// DataService интерфейс сервиса взаимодействия с секретными данными.
type DataService interface {
	Create(ctx context.Context, data smodel.Data) (*smodel.Data, error)
	FindByUUID(ctx context.Context, userUUID, uuid string) (*smodel.Data, error)
	FindByUserUUID(ctx context.Context, userUUID string) ([]*smodel.Data, error)
	Update(ctx context.Context, userUUID, uuid string, value []byte, version time.Time) (*smodel.Data, error)
	Delete(ctx context.Context, userUUID, uuid string) error
}

// Data структура обработчика взаимодействия с секретными данными.
type Data struct {
	service DataService
	logger  log.Loggable
}

// NewData конструктор.
func NewData(service DataService, logger log.Loggable) *Data {
	return &Data{
		service: service,
		logger:  logger,
	}
}

// Create обработчик создания данных.
func (d Data) Create(rd model.DataCreateRequest, w http.ResponseWriter, r *http.Request) (any, int) {
	userUUID, ok := middleware.GetUserUUID(r.Context())
	if !ok {
		return nil, http.StatusUnauthorized
	}

	dt := smodel.Data{
		UserUUID:  userUUID,
		Type:      rd.Type,
		Value:     rd.Value,
		Version:   rd.Version,
		CreatedAt: rd.CreatedAt,
	}
	result, err := d.service.Create(r.Context(), dt)

	if err != nil {
		d.logger.Error("Ошибка при добавлении записи.", err.Error(), dt)
		return nil, http.StatusInternalServerError
	}

	d.logger.Info("Запись успешно добавлена.", result)
	return result, http.StatusOK
}

// Update обработчик обновления данных по uuid.
func (d Data) Update(rd model.DataUpdateRequest, w http.ResponseWriter, r *http.Request) (any, int) {
	userUUID, ok := middleware.GetUserUUID(r.Context())
	if !ok {
		return nil, http.StatusUnauthorized
	}

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		return nil, http.StatusBadRequest
	}

	result, err := d.service.Update(r.Context(), userUUID, uuid, rd.Value, rd.Version)
	if err != nil {
		errCtx := struct {
			UserUUID string
			UUID     string
			Value    []byte
			Version  time.Time
		}{
			UserUUID: userUUID,
			UUID:     uuid,
			Value:    rd.Value,
			Version:  rd.Version,
		}
		d.logger.Error("Ошибка при обновлении записи.", err.Error(), errCtx)
		return nil, http.StatusInternalServerError
	}

	d.logger.Info("Запись успешно обновлена.", result)
	return result, http.StatusOK
}

// Get обработчик получения данных по uuid.
func (d Data) Get(w http.ResponseWriter, r *http.Request) (any, int) {
	userUUID, ok := middleware.GetUserUUID(r.Context())
	if !ok {
		return nil, http.StatusUnauthorized
	}

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		return nil, http.StatusBadRequest
	}

	result, err := d.service.FindByUUID(r.Context(), userUUID, uuid)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, http.StatusNotFound
		}

		errCtx := struct {
			UserUUID string
			UUID     string
		}{
			UserUUID: userUUID,
			UUID:     uuid,
		}

		d.logger.Error("Ошибка при получении записи.", err.Error(), errCtx)
		return nil, http.StatusInternalServerError
	}

	d.logger.Info("Запись успешно получена.")
	return result, http.StatusOK
}

// GetList обработчик получения списка данных.
func (d Data) GetList(w http.ResponseWriter, r *http.Request) (any, int) {
	userUUID, ok := middleware.GetUserUUID(r.Context())
	if !ok {
		return nil, http.StatusUnauthorized
	}

	result, err := d.service.FindByUserUUID(r.Context(), userUUID)
	if err != nil {
		errCtx := struct {
			UserUUID string
		}{
			UserUUID: userUUID,
		}
		d.logger.Error("Ошибка при получении списка записей.", err.Error(), errCtx)
		return nil, http.StatusInternalServerError
	}

	d.logger.Info("Список записей успешно получена.")
	return result, http.StatusOK
}

// Delete обработчик удаления данных по uuid.
func (d Data) Delete(w http.ResponseWriter, r *http.Request) (any, int) {
	userUUID, ok := middleware.GetUserUUID(r.Context())
	if !ok {
		return nil, http.StatusUnauthorized
	}

	uuid := chi.URLParam(r, "uuid")
	if uuid == "" {
		return nil, http.StatusBadRequest
	}

	err := d.service.Delete(r.Context(), userUUID, uuid)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			return nil, http.StatusNotFound
		}

		errCtx := struct {
			UserUUID string
			UUID     string
		}{
			UserUUID: userUUID,
			UUID:     uuid,
		}
		d.logger.Error("Ошибка при удалении записи.", err.Error(), errCtx)
		return nil, http.StatusInternalServerError
	}

	d.logger.Info("Запись успешно удалена.")
	return nil, http.StatusOK
}
