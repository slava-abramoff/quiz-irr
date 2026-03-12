package usecases

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"strings"
)

type AuthServiceProvider interface {
	HashPassword(password string) string
	ComparePassword(hash, password string) bool
	MakeAccessToken(admin *models.Admin) (string, error)
	MakeRefreshToken(admin *models.Admin) (string, error)
	GetPayload(tokenString string) (*dto.TokenPayload, error)
}

type UserServiceProvider interface {
	Create(
		ctx context.Context,
		user dto.CreateUserRequest,
		root bool,
	) (*models.Admin, error)
	GetByEmail(ctx context.Context, email string) (*models.Admin, error)
	GetByID(ctx context.Context, id uint) (*models.Admin, error)
	Update(
		ctx context.Context,
		id uint,
		user dto.UpdateUserRequest,
	) (*models.Admin, error)
	Delete(ctx context.Context, id uint) (*models.Admin, error)
}

type usersCases struct {
	usersSerivce UserServiceProvider
	authService  AuthServiceProvider
}

func NewUsersCases(uS UserServiceProvider, aS AuthServiceProvider) *usersCases {
	return &usersCases{usersSerivce: uS, authService: aS}
}

func (u *usersCases) Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.usersSerivce.GetByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	if !u.authService.ComparePassword(user.Password, data.Password) {
		return nil, errors.New("Invalid password")
	}

	access, err := u.authService.MakeAccessToken(user)
	if err != nil {
		return nil, err
	}

	refresh, err := u.authService.MakeRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		FullName:     user.FullName,
		Email:        user.Email,
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *usersCases) Refresh(ctx context.Context) {}

func (u *usersCases) Create(ctx context.Context, data dto.CreateUserRequest) (*dto.UserResponse, error) {
	data.Password = u.authService.HashPassword(data.Password)

	user, err := u.usersSerivce.Create(ctx, data, false)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		IsRoot:   user.IsRoot,
	}, nil
}

func (u *usersCases) Update(ctx context.Context, id uint, data dto.UpdateUserRequest) (*dto.UserResponse, error) {
	oldUser, err := u.usersSerivce.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if oldUser.IsRoot {
		return nil, fmt.Errorf("Don't permissions for updating root")
	}

	user, err := u.usersSerivce.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		IsRoot:   user.IsRoot,
	}, nil
}

func (u *usersCases) Delete(ctx context.Context, id uint) (*dto.UserResponse, error) {
	oldUser, err := u.usersSerivce.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if oldUser.IsRoot {
		return nil, fmt.Errorf("Don't permissions for delete root")
	}

	user, err := u.usersSerivce.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		IsRoot:   user.IsRoot,
	}, nil
}

func (u *usersCases) GetAuthInfo(r *http.Request) (uint, bool, error) {
	authHeader := r.Header.Get("Authorization")

	const prefix = "Bearer "
	token := strings.TrimPrefix(authHeader, prefix)

	payload, err := u.authService.GetPayload(token)
	if err != nil {
		return 0, false, err
	}

	return payload.ID, payload.IsRoot, nil
}

func (u *usersCases) BootstrapCreate(dto dto.CreateUserRequest) (*models.Admin, error) {
	ctx := context.Background()
	dto.Password = u.authService.HashPassword(dto.Password)
	return u.usersSerivce.Create(ctx, dto, true)
}
