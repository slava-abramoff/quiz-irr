package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/validator"
	"quiz-irr/pkg/httpresponse"

	"github.com/julienschmidt/httprouter"
)

type AuthProvider interface {
	Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type AuthHandlers struct {
	auth AuthProvider
}

func NewAuthHandlers(a AuthProvider) *AuthHandlers {
	return &AuthHandlers{auth: a}
}

func (a *AuthHandlers) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !validator.IsValidEmail(req.Email) {
		httpresponse.ErrorResponse(w, fmt.Errorf("Invalid email").Error(), 422)
		return
	}

	if !validator.IsValidLength(req.Password, 16, 7) {
		httpresponse.ErrorResponse(w, fmt.Errorf("Invalid password").Error(), 422)
		return
	}

	data, err := a.auth.Login(ctx, req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusFound)
		return
	}

	httpresponse.JsonResponse(w, data, 200)
}
