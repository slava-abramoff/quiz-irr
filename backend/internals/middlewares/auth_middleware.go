package middlewares

import (
	"context"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type AuthService interface {
	GetPayload(tokenString string) (*dto.TokenPayload, error)
}

type contextKey string

const AdminContextKey contextKey = "admin"

func NewAuthMiddleware(auth AuthService) func(onlyRoot bool) func(httprouter.Handle) httprouter.Handle {
	return func(onlyRoot bool) func(httprouter.Handle) httprouter.Handle {

		return func(next httprouter.Handle) httprouter.Handle {

			return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					http.Error(w, "missing authorization header", http.StatusUnauthorized)
					return
				}

				const prefix = "Bearer "

				if !strings.HasPrefix(authHeader, prefix) {
					http.Error(w, "invalid authorization header", http.StatusUnauthorized)
					return
				}

				tokenString := strings.TrimPrefix(authHeader, prefix)

				payload, err := auth.GetPayload(tokenString)
				if err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}

				if onlyRoot && !payload.IsRoot {
					http.Error(w, "root access required", http.StatusForbidden)
					return
				}

				ctx := context.WithValue(r.Context(), AdminContextKey, payload)

				next(w, r.WithContext(ctx), ps)
			}
		}
	}
}
