// Package simple содержит вспомогательные функции для работы с http.
package simple

import (
	"encoding/json"
	"net/http"

	"github.com/casnerano/seckeep/pkg/svalid"
)

// TypedHandler функция-адаптер для хендлеров, обеспечивающий типизированные входные данные.
// Инкапсулирует работу с телом запроса, и берет на себя marshal/unmarshal json данных.
// Клиенской хендлер получается тоньще и легче в тестировании.
func TypedHandler[T any](handler func(T, http.ResponseWriter, *http.Request) (any, int)) http.HandlerFunc {
	validator := svalid.New()
	return func(w http.ResponseWriter, r *http.Request) {
		var t T
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			Error(w, http.StatusBadRequest)
			return
		}

		if err := validator.Validate(t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		responseData, statusCode := handler(t, w, r)

		if statusCode > 0 {
			w.WriteHeader(statusCode)
		}

		if responseData != nil {
			b, err := json.Marshal(responseData)
			if err != nil {
				Error(w, http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(b)
		}
	}
}

// Handler функция-адаптер для хендлеров.
func Handler(handler func(http.ResponseWriter, *http.Request) (any, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseData, statusCode := handler(w, r)

		if statusCode > 0 {
			w.WriteHeader(statusCode)
		}

		if responseData != nil {
			b, err := json.Marshal(responseData)
			if err != nil {
				Error(w, http.StatusInternalServerError)
				return
			}
			_, _ = w.Write(b)
		}
	}
}

// Error функция-helper по коду ошибки возвращает статус и текст ошибки.
func Error(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
