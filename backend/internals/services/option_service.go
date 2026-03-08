package services

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
)

type OptionStorageProvider interface {
	Create(
		ctx context.Context,
		questionId uint,
		text string,
		isCorrect bool,
	) (*models.Option, error)
	GetByQuestionID(
		ctx context.Context,
		questionId uint,
	) ([]models.Option, error)
	GetByID(
		ctx context.Context,
		id uint,
	) (*models.Option, error)
	Update(
		ctx context.Context,
		id uint,
		updates map[string]any,
	) (*models.Option, error)
	Delete(
		ctx context.Context,
		id uint,
	) (*models.Option, error)
}

type optionService struct {
	storage OptionStorageProvider
}

func NewOptionService(s OptionStorageProvider) *optionService {
	return &optionService{storage: s}
}

func (o *optionService) Create(ctx context.Context, questionId uint, data dto.CreateOptionRequest) (*models.Option, error) {
	return o.storage.Create(ctx, questionId, data.Text, data.IsCorrect)
}
func (o *optionService) Update(ctx context.Context, id uint, data dto.UpdateOptionRequest) (*models.Option, error) {
	updated := make(map[string]any)

	if data.IsCorrect != nil {
		updated["is_correct"] = *data.IsCorrect
	}

	if data.Text != nil {
		updated["text"] = *data.Text
	}

	return o.storage.Update(ctx, id, updated)
}
func (o *optionService) GetByQuestionID(ctx context.Context, questionId uint) ([]models.Option, error) {
	return o.storage.GetByQuestionID(ctx, questionId)
}
func (o *optionService) Delete(ctx context.Context, id uint) (*models.Option, error) {
	return o.storage.Delete(ctx, id)
}
