// Package data содержит методы для работы с данными из кального хранилища.
// В процессе сохранения и чтения производит необходимую шифрацию и дешифрацию.
package data

//go:generate mockgen -destination=mock/data.go -source=data.go

import (
	"errors"
	"time"

	"github.com/casnerano/seckeep/internal/client/model"
	smodel "github.com/casnerano/seckeep/internal/pkg/model"
)

// Storage интерфейс работы с локальным хранилищем.
type Storage interface {
	Create(storeData *model.StoreData) error
	Read(index int) (*model.StoreData, error)
	GetList() []*model.StoreData
	Update(index int, dataValue []byte, version time.Time) error
	Delete(index int) error
}

// Encryptor интерфейс шифрации и дешифрации.
type Encryptor interface {
	Encrypt(dt model.DataTypeable) ([]byte, error)
	Decrypt(encrypted []byte, dt model.DataTypeable) error
}

// Data структура работы с данными.
type Data struct {
	storage   Storage
	encryptor Encryptor
}

// New конструктор.
func New(storage Storage, encryptor Encryptor) *Data {
	return &Data{
		storage:   storage,
		encryptor: encryptor,
	}
}

// Create метод создает запись.
func (d Data) Create(dt model.DataTypeable) error {
	encrypted, err := d.encryptor.Encrypt(dt)
	if err != nil {
		return err
	}

	sd := &model.StoreData{
		Type:      dt.Type(),
		Value:     encrypted,
		Version:   time.Now(),
		CreatedAt: time.Now(),
	}

	return d.storage.Create(sd)
}

// Read метод читает запись.
func (d Data) Read(index int) (model.DataTypeable, error) {
	storeData, err := d.storage.Read(index)
	if err != nil {
		return nil, err
	}

	var dt model.DataTypeable
	switch storeData.Type {
	case smodel.DataTypeCredential:
		dt = &model.DataCredential{}
	case smodel.DataTypeText:
		dt = &model.DataText{}
	case smodel.DataTypeCard:
		dt = &model.DataCard{}
	case smodel.DataTypeDocument:
		dt = &model.DataDocument{}
	default:
		return nil, errors.New("unknown data type")
	}

	if err = d.encryptor.Decrypt(storeData.Value, dt); err != nil {
		return nil, err
	}

	return dt, nil
}

// GetList метод читает список данных.
func (d Data) GetList() map[int]model.DataTypeable {
	result := make(map[int]model.DataTypeable)
	for index, value := range d.storage.GetList() {
		if value.Deleted {
			continue
		}

		if dt, err := d.Read(index); err == nil {
			result[index] = dt
		}
	}
	return result
}

// Update метод обновляет данные.
func (d Data) Update(index int, dt model.DataTypeable) error {
	encrypted, err := d.encryptor.Encrypt(dt)
	if err != nil {
		return err
	}

	return d.storage.Update(index, encrypted, time.Now())
}

// Delete метод удаляет данные.
func (d Data) Delete(index int) error {
	return d.storage.Delete(index)
}
