package storage

import (
	"context"
	"errors"
	"quiz-irr/internals/storage/models"

	"gorm.io/gorm"
)

type usersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *usersRepo {
	return &usersRepo{db: db}
}

// Создание администратора
func (u *usersRepo) Create(
	ctx context.Context,
	fullName,
	email,
	password string,
	root bool,
) (*models.Admin, error) {

	admin := &models.Admin{
		FullName: fullName,
		Email:    email,
		Password: password,
		IsRoot:   root,
	}

	if err := u.db.WithContext(ctx).Create(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Получение администратора по ID
func (u *usersRepo) GetByID(ctx context.Context, id uint) (*models.Admin, error) {
	var admin models.Admin

	if err := u.db.WithContext(ctx).
		First(&admin, id).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}

// Получение администратора по Email
func (u *usersRepo) GetByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin

	err := u.db.WithContext(ctx).
		Where("email = ?", email).
		First(&admin).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &admin, nil
}

// Обновление администратора
func (u *usersRepo) Update(
	ctx context.Context,
	id uint,
	updated map[string]any,
) (*models.Admin, error) {

	var admin models.Admin

	err := u.db.WithContext(ctx).
		Model(&admin).
		Where("id = ?", id).
		Updates(updated).Error
	if err != nil {
		return nil, err
	}

	err = u.db.WithContext(ctx).
		First(&admin, id).Error
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

// Удаление администратора
func (u *usersRepo) Delete(ctx context.Context, id uint) (*models.Admin, error) {
	var admin models.Admin

	if err := u.db.WithContext(ctx).First(&admin, id).Error; err != nil {
		return nil, err
	}

	if err := u.db.WithContext(ctx).Delete(&admin).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}
