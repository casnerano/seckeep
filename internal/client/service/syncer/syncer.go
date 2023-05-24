// Package syncer дает методы для синхронизации данных локального хранилища с сервером.
package syncer

//go:generate mockgen -destination=mock/syncer.go -source=syncer.go

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-resty/resty/v2"
)

// Основные ошибки при синхронизации.
var (
	// ErrUnauthorized нет авторизации.
	ErrUnauthorized = errors.New("unauthorized")
)

// Storage интерфейс локального хранилища.
type Storage interface {
	Len() int
	OverwriteStore(memStore []*model.StoreData) error
	GetList() []*model.StoreData
}

// Syncer структура синхронизатора.
type Syncer struct {
	serverHealthErr error
	client          *resty.Client
	storage         Storage
	logger          log.Loggable
}

type storeData struct {
	data  *model.StoreData
	index int
}

// New конструктор синхронизатора.
func New(client *resty.Client, storage Storage, logger log.Loggable) *Syncer {
	return &Syncer{
		client:  client,
		storage: storage,
		logger:  logger,
	}
}

// PingServerHealth проверяет связь с сервером.
func (s *Syncer) PingServerHealth() error {
	var response *resty.Response
	response, s.serverHealthErr = s.client.R().Get("/ping")

	if s.serverHealthErr != nil {
		return s.serverHealthErr
	}

	switch response.StatusCode() {
	case http.StatusUnauthorized:
		s.serverHealthErr = ErrUnauthorized
	case http.StatusOK:
		s.serverHealthErr = nil
	default:
		s.serverHealthErr = errors.New(
			http.StatusText(response.StatusCode()),
		)
	}

	return s.serverHealthErr
}

// ServerHealthErr метод возвращает статус связи с сервером.
func (s *Syncer) ServerHealthErr() error {
	return s.serverHealthErr
}

// Run запускает синхронизацию.
func (s *Syncer) Run() error {
	s.logger.Info("Старт синхронизации..")
	s.logger.Info(fmt.Sprintf("Записей в локальном хранилище — %d", s.storage.Len()))

	serverItems, err := s.exportFromServer()
	if err != nil {
		s.logger.Error("Ошибка выгрузки записей из сервера.", err)
		return err
	}

	s.logger.Info(fmt.Sprintf("Записей на сервере — %d", len(serverItems)))

	serverItemsMap := make(map[string]*model.StoreData)
	for k := range serverItems {
		serverItemsMap[serverItems[k].UUID] = serverItems[k]
	}

	// Новые записи на клиенте, необходимые для загрузки на сервер.
	localCreatedItems := make([]*model.StoreData, 0)

	// Записи на клиенте помеченные на удаление, необхомые удалить на сервере.
	localDeletedItems := make([]*storeData, 0)

	// Остальные записи на клиенте, неободимые для синхронизации версий.
	localOtherItemsMap := make(map[string]*storeData)

	// Записи без UUID — созданы локлаьно.
	// Записи с признаком "Deleted" — нужно удалить на сервере.
	for index, sd := range s.storage.GetList() {
		if sd.UUID == "" {
			localCreatedItems = append(localCreatedItems, sd)
		} else if sd.Deleted {
			localDeletedItems = append(localDeletedItems, &storeData{
				data:  sd,
				index: index,
			})
		} else {
			localOtherItemsMap[sd.UUID] = &storeData{
				data:  sd,
				index: index,
			}
		}
	}

	// Загружаем новые клиентские записи на сервер.
	if err = s.loadToServer(localCreatedItems); err != nil {
		s.logger.Error("Ошибка загрузки новых записей на сервер.", err)
		return err
	}

	s.logger.Info(fmt.Sprintf("Загружено новых записей на сервер — %d", len(localCreatedItems)))

	// Отправляем на сервер для удаления записей, которые клиент пометил на удаление.
	if err = s.removeFromServer(localDeletedItems); err != nil {
		s.logger.Error("Ошибка удаления записей из сервера.", err)
		return err
	}

	s.logger.Info(fmt.Sprintf("Удалено записей из сервера — %d", len(localDeletedItems)))

	// Записи на клиенте имеющие более актуальную версию, необхомые загрузить на сервере.
	localUpdatedItems := make([]*model.StoreData, 0)

	for key := range localOtherItemsMap {
		// Если локальная запись есть на сервере.
		if sItem, ok := serverItemsMap[key]; ok {
			// Если серверная версия записи актуальнее локальной.
			if localOtherItemsMap[key].data.Version.After(sItem.Version) {
				localUpdatedItems = append(localUpdatedItems, localOtherItemsMap[key].data)
			}
		}
	}

	// Загрузка обновленных записей на сервер.
	if err = s.updateToServer(localUpdatedItems); err != nil {
		s.logger.Error("Ошибка загрузки обновленных записей на сервер.", err)
		return err
	}

	s.logger.Info(fmt.Sprintf("Обновлено записей на сервере — %d", len(localUpdatedItems)))

	// Выгружаем обновленные данные из сервера.
	actualServerItems, err := s.exportFromServer()
	if err != nil {
		s.logger.Error("Ошибка выгрузки записей из сервера.", err)
		return err
	}

	// Перезаписываем локальное хранилище актуальными данными из сервера.
	if err = s.storage.OverwriteStore(actualServerItems); err != nil {
		s.logger.Error("Ошибка записи актуальных данных из сервера в локальное хранилище.", err)
		return err
	}

	s.logger.Info("Синхронизация успешно завершена.")
	return nil
}

// RunWithStatus метод запускает синхронизацию и выводит статус в stdout.
func (s *Syncer) RunWithStatus() {
	if err := s.Run(); err != nil {
		fmt.Println("Произошла ошибка во время синхронизации.")
	} else {
		fmt.Println("Данные синхронизированы с сервером.")
	}
	fmt.Println()
}

// loadToServer загружает записи на сервер.
func (s *Syncer) loadToServer(localItems []*model.StoreData) error {
	for key := range localItems {
		_, err := s.client.R().
			SetBody(localItems[key]).
			Post("/data")

		if err != nil {
			return err
		}
	}

	return nil
}

// updateToServer загружает записи на сервер.
func (s *Syncer) updateToServer(localItems []*model.StoreData) error {
	for key := range localItems {
		_, err := s.client.R().
			SetBody(struct {
				Value   []byte    `json:"value"`
				Version time.Time `json:"version"`
			}{localItems[key].Value, localItems[key].Version}).
			Put("/data/" + localItems[key].UUID)

		if err != nil {
			return err
		}
	}

	return nil
}

// loadFromServer выгружает записи из сервер.
func (s *Syncer) exportFromServer() ([]*model.StoreData, error) {
	items := make([]*model.StoreData, 0)
	_, err := s.client.R().
		SetResult(&items).
		Get("/data")

	if err != nil {
		return nil, err
	}

	return items, nil
}

// removeFromServer удаляет записи из сервера.
func (s *Syncer) removeFromServer(items []*storeData) error {
	for key := range items {
		_, err := s.client.R().
			Delete("/data/" + items[key].data.UUID)

		if err != nil {
			return err
		}
	}

	return nil
}
