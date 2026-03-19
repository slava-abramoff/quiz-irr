package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/validator"
	"quiz-irr/pkg/apperrors"
	"quiz-irr/pkg/httpresponse"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UsersProvider interface {
	Create(ctx context.Context, data dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, id uint, data dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) (*dto.UserResponse, error)
}

type UserHandlers struct {
	users UsersProvider
}

func NewUserHandlers(u UsersProvider) *UserHandlers {
	return &UserHandlers{users: u}
}

func (u *UserHandlers) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	err := userValidate(req)
	if err != nil {
		httpresponse.ErrorResponse(w, err.Error(), 422)
		return
	}

	user, err := u.users.Create(ctx, req)
	if err != nil {
		code, msg := apperrors.ToHTTP(err)
		httpresponse.ErrorResponse(w, msg, code)
		return
	}

	httpresponse.JsonResponse(w, user, 200)
}

func (u *UserHandlers) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	idStr := ps.ByName("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}
	id := uint(idInt)

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	err = userOptionalValidate(req)
	if err != nil {
		httpresponse.ErrorResponse(w, err.Error(), 422)
		return
	}

	user, err := u.users.Update(ctx, id, req)
	if err != nil {
		code, msg := apperrors.ToHTTP(err)
		httpresponse.ErrorResponse(w, msg, code)
		return
	}

	httpresponse.JsonResponse(w, user, 200)
}

func (u *UserHandlers) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	idStr := ps.ByName("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := uint(idInt)

	user, err := u.users.Delete(ctx, id)
	if err != nil {
		code, msg := apperrors.ToHTTP(err)
		httpresponse.ErrorResponse(w, msg, code)
		return
	}

	httpresponse.JsonResponse(w, user, 200)
}

func userValidate(req dto.CreateUserRequest) error {

	if !validator.IsValidEmail(req.Email) {
		return fmt.Errorf("Invalid email")
	}

	if !validator.IsValidLength(req.FullName, 100, 5) {
		return fmt.Errorf("Invalid FullName")
	}

	if !validator.IsValidLength(req.Password, 16, 8) {
		return fmt.Errorf("Invalid password")
	}

	return nil
}

func userOptionalValidate(req dto.UpdateUserRequest) error {
	if req.Email != "" {
		if !validator.IsValidEmail(req.Email) {
			return fmt.Errorf("Invalid email")
		}
	}

	if req.FullName != "" {
		if !validator.IsValidLength(req.FullName, 100, 5) {
			return fmt.Errorf("Invalid FullName")
		}
	}

	if req.Password != "" {
		if !validator.IsValidLength(req.Password, 16, 8) {
			return fmt.Errorf("Invalid password")
		}
	}

	return nil
}
