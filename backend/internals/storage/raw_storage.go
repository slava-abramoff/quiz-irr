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
	start time.Time,
) (*models.RawSubmission, error) {

	sub := &models.RawSubmission{
		TestID:   testId,
		FullName: fn,
		Email:    email,
		Org:      org,
		StartAt:  &start,
	}

	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *rawRepo) GetByTestID(ctx context.Context, testId uuid.UUID) ([]models.RawSubmission, error) {
	var subs []models.RawSubmission

	err := r.db.WithContext(ctx).
		Where("test_id = ?", testId).
		Find(&subs).Error

	if err != nil {
		return nil, err
	}

	return subs, nil
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
