package handler

import (
	"net/http"

	"github.com/casnerano/seckeep/pkg/log"
)

// Service структура обработчика взаимодействия с сервисом (сервером).
type Service struct {
	logger log.Loggable
}

// NewService конструктор.
func NewService(logger log.Loggable) *Service {
	return &Service{logger: logger}
}

// Ping обработчик проверки.
func (s *Service) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
