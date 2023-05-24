package http

import (
	"github.com/casnerano/seckeep/internal/server/http/handler"
	"github.com/casnerano/seckeep/internal/server/http/middleware"
	"github.com/casnerano/seckeep/internal/server/service/account"
	"github.com/casnerano/seckeep/internal/server/service/data"
	"github.com/casnerano/seckeep/pkg/http/simple"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// Router структура роутера.
type Router struct {
	chiRouter *chi.Mux
	logger    *log.Logger
	secret    string
}

// NewRouter конструктор.
func NewRouter(logger *log.Logger, secret string) *Router {
	chiRouter := chi.NewRouter()

	chiRouter.Use(chiMiddleware.RequestID)
	chiRouter.Use(chiMiddleware.Recoverer)
	chiRouter.Use(middleware.JSONContentType)

	return &Router{
		chiRouter: chiRouter,
		logger:    logger,
		secret:    secret,
	}
}

// InitServiceHandler метод инициализации роутов для обработчиков сервиса (сервера).
func (router *Router) InitServiceHandler() {
	h := handler.NewService(router.logger)
	router.chiRouter.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator(router.secret))
		r.Get("/api/ping", h.Ping())
	})
}

// InitAccountHandler метод инициализации роутов для обработчиков аккаунта пользователя.
func (router *Router) InitAccountHandler(service *account.Account) {
	h := handler.NewAccount(service, router.logger)
	router.chiRouter.Group(func(r chi.Router) {
		r.Post("/api/user/register", simple.TypedHandler(h.SignUp))
		r.Post("/api/user/login", simple.TypedHandler(h.SignIn))
	})
}

// InitDataHandler метод инициализации роутов для обработчиков взаимодействия с секретными данными.
func (router *Router) InitDataHandler(service *data.Data) {
	h := handler.NewData(service, router.logger)
	router.chiRouter.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator(router.secret))
		r.Post("/api/data", simple.TypedHandler(h.Create))
		r.Get("/api/data", simple.Handler(h.GetList))
		r.Put("/api/data/{uuid}", simple.TypedHandler(h.Update))
		r.Get("/api/data/{uuid}", simple.Handler(h.Get))
		r.Delete("/api/data/{uuid}", simple.Handler(h.Delete))
	})
}

// GetChiMux возвращает дефолтный роутер (chi.Mux).
func (router *Router) GetChiMux() *chi.Mux {
	return router.chiRouter
}
