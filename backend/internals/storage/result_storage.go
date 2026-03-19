package storage

import (
	"context"
	"quiz-irr/internals/handlers/dto"
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
	totalScore int,
	isOnTime bool,
) (*models.TestResult, error) {
	result := &models.TestResult{
		TestID:     testId,
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

func (r *resultRepo) CreateMany(ctx context.Context, results []dto.CreateResult) error {
	if len(results) == 0 {
		return nil
	}

	records := make([]models.TestResult, 0, len(results))

	for _, res := range results {
		records = append(records, models.TestResult{
			TestID:     res.TestID,
			FullName:   res.FullName,
			Email:      res.Email,
			Org:        res.Org,
			Duration:   res.Duration,
			TotalScore: res.TotalScore,
			IsOnTime:   res.IsOnTime,
		})
	}

	if err := r.db.WithContext(ctx).Create(&records).Error; err != nil {
		return err
	}

	return nil
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

func (r *resultRepo) GetAllByTestID(
	ctx context.Context,
	testId uuid.UUID,
) ([]models.TestResult, error) {

	var results []models.TestResult

	err := r.db.WithContext(ctx).
		Where("test_id = ?", testId).
		Order("id DESC").
		Find(&results).
		Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *resultRepo) GetByTestID(
	ctx context.Context,
	testId uuid.UUID,
	skip, take uint,
) ([]models.TestResult, uint, error) {

	var (
		results []models.TestResult
		count   int64
	)

	db := r.db.WithContext(ctx).Model(&models.TestResult{}).
		Where("test_id = ?", testId).Order("total_score DESC")

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Offset(int(skip)).
		Limit(int(take)).
		Order("id DESC").
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, uint(count), nil
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

func (r *resultRepo) DeleteByTestID(
	ctx context.Context,
	testId uuid.UUID,
) error {

	return r.db.WithContext(ctx).
		Where("test_id = ?", testId).
		Delete(&models.TestResult{}).
		Error
}
