// Package client выступает в качестве основы для клиентской части приложения.
// Выполняет инициализацию, конфигурации зависимостей и управление жизненым циклом.
package client

import (
	"fmt"
	"os"

	"github.com/casnerano/seckeep/internal/client/command"
	"github.com/casnerano/seckeep/internal/client/config"
	"github.com/casnerano/seckeep/internal/client/service/storage"
	"github.com/casnerano/seckeep/pkg/config/yaml"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/casnerano/seckeep/pkg/log/handler"
	"github.com/casnerano/seckeep/pkg/log/handler/formatter"
)

// App структура приложения.
type App struct {
	config      *config.Config
	logger      *log.Logger
	dataStorage *storage.Storage
	rootCmd     *command.Root
}

// NewApp конструктор.
func NewApp() (*App, error) {
	var err error
	app := &App{}

	// Инициализация логгера.
	app.logger = log.New("client")

	app.logger.AddHandler(
		handler.NewStdOut(
			formatter.NewText(),
			log.LevelError,
			true,
		),
	)

	// Инициализация конфигурации.
	app.config = &config.Config{}
	if err = yaml.LoadFromFile(config.FileName, app.config); err != nil {
		app.logger.Emergency("Не удалось прочитать файл конфигурации.", err)
		return nil, err
	}

	// Инициализация локального хранилища.
	app.dataStorage, err = storage.New(storage.DefaultFileName)
	if err != nil {
		return nil, err
	}

	// Инициализация рутовой команды.
	app.rootCmd = command.NewRoot(&command.RootCommandContext{
		Config:      app.config,
		Logger:      app.logger,
		DataStorage: app.dataStorage,
	})

	return app, nil
}

// Run метод запуска приложения.
func (a *App) Run() {
	if err := a.rootCmd.Execute(); err != nil {
		a.logger.Emergency("Ошибка запуска клиента.", err)
		os.Exit(1)
	}
}

// Shutdown метод завершения приложения.
func (a *App) Shutdown() {
	if err := a.dataStorage.Close(); err != nil {
		a.logger.Error("Не удалось закрыть локальное хранилище.", err.Error())
	}

	if err := a.logger.Close(); err != nil {
		fmt.Println("Не удалось завершить логгер.", err.Error())
	}
}
