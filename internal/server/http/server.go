package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-chi/chi/v5"
)

// Server структура сервера.
type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

// NewServer конструктор.
func NewServer(address string, router *chi.Mux, logger *log.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: router,
		},
		logger: logger,
	}
}

// Start метод стартует сервер и ожидает завершения контекста.
func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Alert(
				fmt.Sprintf("Ошибка запуска сервера по адресу %s", s.httpServer.Addr),
				err,
			)
			os.Exit(1)
		}
	}()

	s.logger.Info(fmt.Sprintf("Сервер запущен по адресу %s", s.httpServer.Addr))

	<-ctx.Done()
	s.logger.Info("Приостановление обслуживания новых http-запросов..")

	if err := s.Shutdown(); err != nil {
		return err
	}

	s.logger.Info("Сервер завершен.")
	return nil
}

// Shutdown метод завершает все необработанные обработчики, и останавливает сервер.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
