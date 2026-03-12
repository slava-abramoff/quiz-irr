package storage

import (
	"context"
	"quiz-irr/internals/storage/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type testsRepo struct {
	db *gorm.DB
}

func NewTestsRepo(db *gorm.DB) *testsRepo {
	return &testsRepo{db: db}
}

func (t *testsRepo) Create(
	ctx context.Context,
	authorId uint,
	title string,
	desc string,
	isActive bool,
	startAt time.Time,
	endAt time.Time,
) (*models.Test, error) {

	test := &models.Test{
		AuthorID: authorId,
		Title:    title,
		Desc:     desc,
		IsActive: isActive,
		StartAt:  &startAt,
		EndAt:    &endAt,
	}

	if err := t.db.WithContext(ctx).Create(test).Error; err != nil {
		return nil, err
	}

	return test, nil
}

func (t *testsRepo) FindMany(ctx context.Context, skip, take uint) ([]models.Test, uint, error) {
	var tests []models.Test
	var total int64

	db := t.db.WithContext(ctx)

	if err := db.Model(&models.Test{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Preload("Author").
		Limit(int(take)).
		Offset(int(skip)).
		Order("created_at DESC").
		Find(&tests).Error; err != nil {
		return nil, 0, err
	}

	return tests, uint(total), nil
}

func (t *testsRepo) GetByAuthor(ctx context.Context, authorId uint) ([]models.Test, error) {
	var tests []models.Test

	if err := t.db.
		WithContext(ctx).
		Where("author_id = ?", authorId).
		Preload("Author").
		Order("created_at DESC").
		Find(&tests).Error; err != nil {
		return nil, err
	}

	return tests, nil
}

func (t *testsRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Test, error) {
	var test models.Test

	if err := t.db.
		WithContext(ctx).
		Preload("Author").
		Where("id = ?", id).
		First(&test).Error; err != nil {
		return nil, err
	}

	return &test, nil
}

func (t *testsRepo) Update(
	ctx context.Context,
	id uuid.UUID,
	updated map[string]any,
) (*models.Test, error) {

	db := t.db.WithContext(ctx)

	var test models.Test

	if err := db.Model(&test).
		Where("id = ?", id).
		Updates(updated).Error; err != nil {
		return nil, err
	}

	if err := db.
		Preload("Author").
		First(&test, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &test, nil
}

func (t *testsRepo) Delete(ctx context.Context, id uuid.UUID) (*models.Test, error) {
	db := t.db.WithContext(ctx)

	var test models.Test

	if err := db.
		Preload("Author").
		First(&test, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := db.Delete(&test).Error; err != nil {
		return nil, err
	}

	return &test, nil
}
