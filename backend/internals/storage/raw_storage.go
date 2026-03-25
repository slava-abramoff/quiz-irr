package storage

import (
	"context"
	"quiz-irr/internals/storage/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type rawRepo struct {
	db *gorm.DB
}

func NewRawRepo(db *gorm.DB) *rawRepo {
	return &rawRepo{db: db}
}

func (r *rawRepo) Create(
	ctx context.Context,
	testId uuid.UUID,
	fn, email, org string,
	birthYear uint,
	start time.Time,
) (*models.RawSubmission, error) {

	sub := &models.RawSubmission{
		TestID:    testId,
		FullName:  fn,
		Email:     email,
		Org:       org,
		BirthYear: &birthYear,
		StartAt:   &start,
	}

	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *rawRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error) {
	var sub models.RawSubmission

	err := r.db.WithContext(ctx).
		Preload("Test").
		Where("id = ?", id).
		First(&sub).Error

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *rawRepo) GetAllByTest(
	ctx context.Context,
	testId uuid.UUID,
) ([]models.RawSubmission, error) {
	var subs []models.RawSubmission

	err := r.db.WithContext(ctx).
		Preload("Test").
		Where("test_id = ?", testId).
		Order("created_at DESC").
		Find(&subs).
		Error

	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (r *rawRepo) GetByTestID(ctx context.Context, testId uuid.UUID, skip, take uint) ([]models.RawSubmission, uint, error) {
	var (
		subs  []models.RawSubmission
		count int64
	)

	db := r.db.WithContext(ctx).Model(&models.RawSubmission{}).
		Where("test_id = ?", testId)

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Offset(int(skip)).
		Limit(int(take)).
		Order("created_at DESC").
		Find(&subs).Error; err != nil {
		return nil, 0, err
	}

	return subs, uint(count), nil
}

func (r *rawRepo) SavePayload(
	ctx context.Context,
	id uuid.UUID,
	payload datatypes.JSON,
) (*models.RawSubmission, error) {

	var sub models.RawSubmission

	res := r.db.WithContext(ctx).
		Model(&models.RawSubmission{}).
		Where("id = ?", id).
		Update("answers_payload", payload).
		Update("end_at", time.Now()).
		Update("status", "end")

	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err := r.db.WithContext(ctx).
		First(&sub, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *rawRepo) Load(
	ctx context.Context,
	id uuid.UUID,
) (datatypes.JSON, error) {

	var sub models.RawSubmission

	err := r.db.WithContext(ctx).
		Select("answers_payload").
		First(&sub, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return sub.AnswersPayload, nil
}

func (r *rawRepo) Update(ctx context.Context, id uuid.UUID, updated map[string]any) (*models.RawSubmission, error) {
	var sub models.RawSubmission

	err := r.db.WithContext(ctx).
		Model(&sub).
		Where("id = ?", id).
		Updates(updated).
		Error

	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).
		First(&sub, "id = ?", id).
		Error

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *rawRepo) Delete(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error) {
	var sub models.RawSubmission

	if err := r.db.WithContext(ctx).
		First(&sub, "id = ?", id).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).
		Delete(&sub).
		Error; err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *rawRepo) DeleteAll(ctx context.Context, testId uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("test_id = ?", testId).
		Delete(&models.RawSubmission{}).Error
}
