package storage

import (
	"context"
	"quiz-irr/internals/storage/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type questionsRepo struct {
	db *gorm.DB
}

func NewQuestionsRepo(db *gorm.DB) *questionsRepo {
	return &questionsRepo{db: db}
}

func (q *questionsRepo) Create(
	ctx context.Context,
	testId uuid.UUID,
	text string,
	questionType string,
	points int,
) (*models.Question, error) {

	question := &models.Question{
		TestID: testId.String(),
		Text:   text,
		Type:   questionType,
		Points: points,
	}

	if err := q.db.WithContext(ctx).Create(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (q *questionsRepo) GetByID(
	ctx context.Context,
	id uint,
) (*models.Question, error) {
	var question models.Question

	err := q.db.WithContext(ctx).
		First(&question, id).Error

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (q *questionsRepo) GetByTestID(
	ctx context.Context,
	testId uuid.UUID,
) ([]models.Question, error) {

	var questions []models.Question

	if err := q.db.WithContext(ctx).
		Where("test_id = ?", testId).
		Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func (q *questionsRepo) Update(
	ctx context.Context,
	id uint,
	updated map[string]any,
) (*models.Question, error) {

	var question models.Question

	if err := q.db.WithContext(ctx).
		Model(&question).
		Where("id = ?", id).
		Updates(updated).Error; err != nil {
		return nil, err
	}

	if err := q.db.WithContext(ctx).
		First(&question, id).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (q *questionsRepo) Delete(
	ctx context.Context,
	id uint,
) (*models.Question, error) {

	var question models.Question

	if err := q.db.WithContext(ctx).
		First(&question, id).Error; err != nil {
		return nil, err
	}

	if err := q.db.WithContext(ctx).
		Delete(&question).Error; err != nil {
		return nil, err
	}

	return &question, nil
}
