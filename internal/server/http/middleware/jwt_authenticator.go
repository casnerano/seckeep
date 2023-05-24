package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/casnerano/seckeep/pkg/jwtoken"
)

type ctxUserUUIDType string

// CtxUserUUIDKey ключ параметра контекста для UUID пользователя.
const CtxUserUUIDKey ctxUserUUIDType = "user_uuid"

// JWTAuthenticator middleware выполняет аутентификацию по JWT токену из заголовка запроса.
func JWTAuthenticator(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			jwt := jwtoken.New()
			payload, err := jwt.Verify(tokenString, []byte(secret))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserUUIDKey, payload.UUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserUUID функция возвращает UUID пользователя из заданного контекста.
func GetUserUUID(ctx context.Context) (string, bool) {
	uuid, ok := ctx.Value(CtxUserUUIDKey).(string)
	if !ok {
		return "", false
	}
	return uuid, true
}
