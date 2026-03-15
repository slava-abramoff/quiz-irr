package services

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"

	"github.com/google/uuid"
)

type QuestionStorageProvider interface {
	Create(
		ctx context.Context,
		testId uuid.UUID,
		text string,
		questionType string,
		points int,
	) (*models.Question, error)
	Update(
		ctx context.Context,
		id uint,
		updated map[string]any,
	) (*models.Question, error)
	GetByID(
		ctx context.Context,
		id uint,
	) (*models.Question, error)
	GetByTestID(
		ctx context.Context,
		testId uuid.UUID,
	) ([]models.Question, error)
	Delete(
		ctx context.Context,
		id uint,
	) (*models.Question, error)
}

type questionService struct {
	storage QuestionStorageProvider
}

func NewQuestionService(s QuestionStorageProvider) *questionService {
	return &questionService{storage: s}
}

func (q *questionService) Create(ctx context.Context, testId uuid.UUID, data dto.CreateQuestionRequest) (*models.Question, error) {
	return q.storage.Create(
		ctx,
		testId,
		data.Text,
		data.Type,
		data.Points,
	)
}

func (q *questionService) GetByID(ctx context.Context, id uint) (*models.Question, error) {
	return q.storage.GetByID(ctx, id)
}

func (q *questionService) Update(ctx context.Context, id uint, data dto.UpdateQuestionRequest) (*models.Question, error) {
	updated := make(map[string]any)

	if data.Points != nil {
		updated["points"] = *data.Points
	}

	if data.Text != nil {
		updated["text"] = *data.Text
	}

	if data.Type != nil {
		updated["type"] = q.SetupType(*data.Type)
	}

	return q.storage.Update(ctx, id, updated)
}

func (q *questionService) GetByTestID(ctx context.Context, testId uuid.UUID) ([]models.Question, error) {
	return q.storage.GetByTestID(ctx, testId)
}

func (q *questionService) Delete(ctx context.Context, id uint) (*models.Question, error) {
	return q.storage.Delete(ctx, id)
}

func (q *questionService) SetupType(t string) string {
	switch t {
	case "single":
		return "single"
	case "text":
		return "text"
	default:
		return "multiple"
	}
}
