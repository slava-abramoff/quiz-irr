package services

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
)

type UsersStorageProvider interface {
	Create(
		ctx context.Context,
		fullName,
		email,
		password string,
		root bool,
	) (*models.Admin, error)
	GetByID(ctx context.Context, id uint) (*models.Admin, error)
	GetByEmail(ctx context.Context, email string) (*models.Admin, error)
	Update(
		ctx context.Context,
		id uint,
		updated map[string]any,
	) (*models.Admin, error)
	Delete(ctx context.Context, id uint) (*models.Admin, error)
}

type userService struct {
	uRepo UsersStorageProvider
}

func NewUsersService(uRepo UsersStorageProvider) *userService {
	return &userService{uRepo: uRepo}
}

func (u *userService) Create(
	ctx context.Context,
	user dto.CreateUserRequest,
	root bool,
) (*models.Admin, error) {
	return u.uRepo.Create(
		ctx,
		user.FullName,
		user.Email,
		user.Password,
		root,
	)
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	return u.uRepo.GetByEmail(ctx, email)
}

func (u *userService) GetByID(ctx context.Context, id uint) (*models.Admin, error) {
	return u.uRepo.GetByID(ctx, id)
}

func (u *userService) Update(
	ctx context.Context,
	id uint,
	user dto.UpdateUserRequest,
) (*models.Admin, error) {
	updated := make(map[string]any)

	if user.FullName != "" {
		updated["full_name"] = user.FullName
	}

	if user.Email != "" {
		updated["email"] = user.Email
	}

	if user.Password != "" {
		updated["password"] = user.Password
	}

	return u.uRepo.Update(ctx, id, updated)
}

func (u *userService) Delete(ctx context.Context, id uint) (*models.Admin, error) {
	return u.uRepo.Delete(ctx, id)
}
