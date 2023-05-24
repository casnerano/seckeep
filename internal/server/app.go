// Package server выступает в качестве основы для серверной части приложения.
// Выполняет инициализацию, конфигурации зависимостей и управление жизненым циклом.
package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/casnerano/seckeep/internal/server/config"
	"github.com/casnerano/seckeep/internal/server/http"
	"github.com/casnerano/seckeep/internal/server/repository/pgsql"
	"github.com/casnerano/seckeep/internal/server/service/account"
	"github.com/casnerano/seckeep/internal/server/service/data"
	"github.com/casnerano/seckeep/pkg/config/yaml"
	"github.com/casnerano/seckeep/pkg/jwtoken"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/casnerano/seckeep/pkg/log/handler"
	"github.com/casnerano/seckeep/pkg/log/handler/formatter"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App структура приложения.
type App struct {
	config  *config.Config
	logger  *log.Logger
	server  *http.Server
	pgxpool *pgxpool.Pool
}

// NewApp конструктор.
func NewApp() (*App, error) {
	var err error
	app := App{}

	// Инициализация логгера.
	app.logger = log.New("server")

	app.logger.AddHandler(
		handler.NewStdOut(
			formatter.NewText(),
			log.LevelDebug,
			true,
		),
	)

	// Инициализация конфигурации.
	app.config = &config.Config{}
	if err = yaml.LoadFromFile(config.FileName, app.config); err != nil {
		app.logger.Emergency("Не удалось прочитать файл конфигурации.", err)
		return nil, err
	}

	// Подключение к базе данных.
	app.pgxpool, err = pgxpool.New(context.Background(), app.config.Database.DSN)
	if err != nil {
		app.logger.Emergency("Не удалось подключиться к базе данных.", err)
		return nil, err
	}

	// Инициализация зависимостей.

	userRepository := pgsql.NewUserRepository(app.pgxpool)
	dataRepository := pgsql.NewDataRepository(app.pgxpool)

	accountService := account.New(
		userRepository,
		jwtoken.New(),
		app.config.App.Authenticator.Secret,
	)

	router := http.NewRouter(app.logger, app.config.App.Authenticator.Secret)
	router.InitServiceHandler()
	router.InitAccountHandler(accountService)
	router.InitDataHandler(data.New(dataRepository))

	app.server = http.NewServer(app.config.Server.Addr, router.GetChiMux(), app.logger)

	return &app, nil
}

// Run метод запуска приложения.
func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := a.server.Start(ctx); err != nil {
		a.logger.Emergency("Ошибка запуска сервера.", err)
		os.Exit(1)
	}
}

// Shutdown метод завершения приложения.
func (a *App) Shutdown() {
	_ = a.server.Shutdown()
	_ = a.logger.Close()
	a.pgxpool.Close()
}
