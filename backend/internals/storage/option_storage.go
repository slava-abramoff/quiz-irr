package storage

import (
	"context"
	"quiz-irr/internals/storage/models"

	"gorm.io/gorm"
)

type optionsRepo struct {
	db *gorm.DB
}

func NewOptionsRepo(db *gorm.DB) *optionsRepo {
	return &optionsRepo{db: db}
}

func (o *optionsRepo) Create(
	ctx context.Context,
	questionId uint,
	text string,
	isCorrect bool,
) (*models.Option, error) {

	option := &models.Option{
		QuestionID: questionId,
		Text:       text,
		IsCorrect:  isCorrect,
	}

	if err := o.db.WithContext(ctx).Create(option).Error; err != nil {
		return nil, err
	}

	return option, nil
}

func (o *optionsRepo) GetByQuestionID(
	ctx context.Context,
	questionId uint,
) ([]models.Option, error) {
	var options []models.Option

	err := o.db.
		WithContext(ctx).
		Where("question_id = ?", questionId).
		Find(&options).
		Error

	if err != nil {
		return nil, err
	}

	return options, nil
}

func (o *optionsRepo) GetByID(
	ctx context.Context,
	id uint,
) (*models.Option, error) {
	var option models.Option

	err := o.db.
		WithContext(ctx).
		First(&option, id).
		Error

	if err != nil {
		return nil, err
	}

	return &option, nil
}

func (o *optionsRepo) Update(
	ctx context.Context,
	id uint,
	updates map[string]any,
) (*models.Option, error) {
	var option models.Option

	if err := o.db.WithContext(ctx).First(&option, id).Error; err != nil {
		return nil, err
	}

	if err := o.db.WithContext(ctx).Model(&option).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &option, nil
}

func (o *optionsRepo) Delete(
	ctx context.Context,
	id uint,
) (*models.Option, error) {
	var option models.Option

	if err := o.db.WithContext(ctx).First(&option, id).Error; err != nil {
		return nil, err
	}

	if err := o.db.WithContext(ctx).Delete(&option).Error; err != nil {
		return nil, err
	}

	return &option, nil
}
