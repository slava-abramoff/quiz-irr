package storage

import (
	"context"
	"quiz-irr/internals/storage/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resultRepo struct {
	db *gorm.DB
}

func NewResultRepo(db *gorm.DB) *resultRepo {
	return &resultRepo{db: db}
}

func (r *resultRepo) Create(
	ctx context.Context,
	testId uuid.UUID,
	fullName string,
	email string,
	org string,
	duration uint,
	totalScore uint,
	isOnTime bool,
) (*models.TestResult, error) {
	result := &models.TestResult{
		TestID:     testId.String(),
		FullName:   fullName,
		Email:      email,
		Org:        org,
		Duration:   duration,
		TotalScore: totalScore,
		IsOnTime:   isOnTime,
	}

	if err := r.db.WithContext(ctx).Create(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *resultRepo) GetByID(
	ctx context.Context,
	id uint,
) (*models.TestResult, error) {
	result := &models.TestResult{}

	if err := r.db.WithContext(ctx).First(result, id).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *resultRepo) GetByTestID(
	ctx context.Context,
	testId uuid.UUID,
) ([]models.TestResult, error) {
	var results []models.TestResult

	if err := r.db.WithContext(ctx).
		Where("test_id = ?", testId.String()).
		Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (r *resultRepo) Update(
	ctx context.Context,
	id uint,
	updates map[string]any,
) (*models.TestResult, error) {
	result := &models.TestResult{}

	if err := r.db.WithContext(ctx).First(result, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(result).Updates(updates).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *resultRepo) Delete(
	ctx context.Context,
	id uint,
) (*models.TestResult, error) {
	result := &models.TestResult{}

	if err := r.db.WithContext(ctx).First(result, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Delete(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
