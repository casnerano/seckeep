package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"
)

// Server структура сервера.
type Server struct {
	httpServer  *http.Server
	enableHTTPS bool
	logger      *log.Logger
}

// NewServer конструктор.
func NewServer(address string, enableHTTPS bool, router *chi.Mux, logger *log.Logger) *Server {
	server := Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: router,
		},
		enableHTTPS: enableHTTPS,
		logger:      logger,
	}

	if enableHTTPS {
		certManager := &autocert.Manager{
			Cache:      autocert.DirCache("./var"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(),
		}
		server.httpServer.TLSConfig = certManager.TLSConfig()
	}

	return &server
}

// Start метод стартует сервер и ожидает завершения контекста.
func (s *Server) Start(ctx context.Context) error {
	go func() {
		if s.enableHTTPS {
			if err := s.httpServer.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				s.logger.Alert(
					fmt.Sprintf("Ошибка запуска сервера по адресу %s", s.httpServer.Addr),
					err,
				)
			}
		} else {
			if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.logger.Alert(
					fmt.Sprintf("Ошибка запуска сервера по адресу %s", s.httpServer.Addr),
					err,
				)
			}
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
