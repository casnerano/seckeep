// Package storage дает низкоуровневые методы работы с локальным хранилищем.
package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/casnerano/seckeep/internal/client/model"
)

const (
	// DefaultFileName дефолтный путь к файлу храненого хранилища.
	DefaultFileName = "./cmd/client/var/store/data.registry"
)

// Основные ошибки при работе с локальным хранилищем.
var (
	// ErrOutOfRangeStore вышел за пределы индекса данных.
	ErrOutOfRangeStore = errors.New("out of range store")
)

// Storage структура работы с локальными хранилищем.
type Storage struct {
	fileStore *os.File
	memStore  []*model.StoreData
	bufStore  *bufio.ReadWriter
}

// New конструктор.
func New(fName string) (*Storage, error) {
	file, err := os.OpenFile(fName, os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(file)

	storage := &Storage{
		fileStore: file,
		memStore:  make([]*model.StoreData, 0),
		bufStore:  bufio.NewReadWriter(reader, writer),
	}

	if err = storage.restore(); err != nil {
		return nil, err
	}

	return storage, nil
}

// Len метод выводит кол-во записей.
func (s *Storage) Len() int {
	return len(s.memStore)
}

// OverwriteStore метод перезаписывает все данные из заданного слайса.
func (s *Storage) OverwriteStore(memStore []*model.StoreData) error {
	s.memStore = memStore
	return s.ClearFlush()
}

// Create метод создает запись.
func (s *Storage) Create(storeData *model.StoreData) error {
	if err := json.NewEncoder(s.bufStore).Encode(storeData); err != nil {
		return err
	}

	if err := s.bufStore.Flush(); err != nil {
		return err
	}

	s.memStore = append(s.memStore, storeData)
	return nil
}

// Read метод читает запись по индексу.
func (s *Storage) Read(index int) (*model.StoreData, error) {
	if index >= 0 && index < len(s.memStore) {
		return s.memStore[index], nil
	}
	return nil, ErrOutOfRangeStore
}

// GetList метод возвращает слайс со всеми данными.
func (s *Storage) GetList() []*model.StoreData {
	return s.memStore
}

// Update метод обновляет запись по индексу.
func (s *Storage) Update(index int, dataValue []byte, version time.Time) error {
	savedData := *s.memStore[index]
	s.memStore[index].Value = dataValue
	s.memStore[index].Version = version
	if err := s.ClearFlush(); err != nil {
		s.memStore[index].Value = savedData.Value
		s.memStore[index].Version = savedData.Version
	}
	return nil
}

// Delete метод удаляет запись по индексу.
func (s *Storage) Delete(index int) error {
	s.memStore[index].Deleted = true
	if err := s.ClearFlush(); err != nil {
		s.memStore[index].Deleted = false
	}
	return nil
}

// ClearFlush метод очищает хранилище, и заново записывает все данные из памяти.
func (s *Storage) ClearFlush() error {
	if _, err := s.fileStore.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := s.fileStore.Truncate(0); err != nil {
		return err
	}

	for _, storeData := range s.memStore {
		if err := s.Create(storeData); err != nil {
			return nil
		}
	}

	return s.bufStore.Flush()
}

// Close метод закрывает файл хранилища.
func (s *Storage) Close() error {
	return s.fileStore.Close()
}

// restore метод восстанавливает в память данные из файла хранилища.
func (s *Storage) restore() error {
	for {
		line, err := s.bufStore.ReadSlice('\n')
		if err == io.EOF {
			break
		}

		storeData := &model.StoreData{}
		err = json.Unmarshal(line, storeData)

		if err != nil {
			continue
		}

		s.memStore = append(s.memStore, storeData)
	}

	return nil
}
