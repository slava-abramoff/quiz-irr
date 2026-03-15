package usecases

import (
	"context"
	"errors"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"time"

	"github.com/google/uuid"
)

type TestInfoProvider interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Test, error)
}

type QuestionInfoProvider interface {
	GetByTestID(ctx context.Context, testId uuid.UUID) ([]models.Question, error)
}

type OptionInfoProvider interface {
	GetByQuestionID(ctx context.Context, questionId uint) ([]models.Option, error)
}

type RawDataServiceProvider interface {
	Load(ctx context.Context, id uuid.UUID) (*dto.SendUserAnswersRequest, error)
	SavePayload(ctx context.Context, id uuid.UUID, data dto.SendUserAnswersRequest) (*models.RawSubmission, error)
	Create(ctx context.Context, testId uuid.UUID, f, e, o string, start time.Time) (*models.RawSubmission, error)
}

type examCases struct {
	testService     TestInfoProvider
	optionService   OptionInfoProvider
	questionService QuestionInfoProvider
	rawDataService  RawDataServiceProvider
}

func NewExamCases(t TestInfoProvider, r RawDataServiceProvider, o OptionInfoProvider, q QuestionInfoProvider) *examCases {
	return &examCases{
		testService:     t,
		rawDataService:  r,
		optionService:   o,
		questionService: q,
	}
}

// Получить информацию о тесте
func (e *examCases) GetTestInfo(ctx context.Context, id uuid.UUID) (*dto.TestCustomerResponse, error) {
	test, err := e.testService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !test.IsActive {
		return nil, errors.New("Not found test")
	}

	return &dto.TestCustomerResponse{
		Title:   test.Title,
		Desc:    test.Desc,
		StartAt: test.StartAt.Format("2006-01-02 15:04:05"),
		EndAt:   test.EndAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (e *examCases) StartTest(ctx context.Context, id uuid.UUID, data dto.StartExamRequest) (*dto.StartExamResponse, error) {
	// TODO: опрашивать кэш
	questions, err := e.questionService.GetByTestID(ctx, id)
	if err != nil {
		return nil, err
	}

	questionDtos := make([]dto.ExamQuestion, 0, 100)

	for _, q := range questions {
		questionDto := dto.ExamQuestion{
			ID:   q.ID,
			Text: q.Text,
			Type: q.Type,
		}

		optionDtos := make([]dto.ExamOption, 0, 10)

		options, err := e.optionService.GetByQuestionID(ctx, q.ID)
		if err != nil {
			return nil, err
		}

		for _, o := range options {
			optionDto := dto.ExamOption{
				ID:   o.ID,
				Text: o.Text,
			}

			optionDtos = append(optionDtos, optionDto)
		}

		questionDto.Options = optionDtos
		questionDtos = append(questionDtos, questionDto)
	}

	startAt := time.Now()

	raw, err := e.rawDataService.Create(
		ctx,
		id,
		data.FullName,
		data.Email,
		data.Org,
		startAt,
	)
	if err != nil {
		return nil, err
	}

	return &dto.StartExamResponse{
		DataID:    raw.ID.String(),
		Questions: questionDtos,
	}, nil
}

func (e *examCases) SaveAnswers(ctx context.Context, rawId uuid.UUID, data dto.SendUserAnswersRequest) (*dto.SendUserAnswersResponse, error) {
	_, err := e.rawDataService.SavePayload(ctx, rawId, data)
	if err != nil {
		return nil, err
	}

	return &dto.SendUserAnswersResponse{
		Message: "Ожидайте результатов",
	}, nil
}

func (e *examCases) isTimeToStart(startAt, endAt time.Time) bool {
	now := time.Now()
	return now.After(startAt) && now.Before(endAt)
}
