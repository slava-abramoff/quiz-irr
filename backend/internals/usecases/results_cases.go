package usecases

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"

	"github.com/google/uuid"
)

type ResultServiceProvider interface {
	GetAll(
		ctx context.Context,
		testId uuid.UUID,
	) ([]models.TestResult, error)
	GetByTestID(
		ctx context.Context,
		testId uuid.UUID,
		skip, take uint,
	) ([]models.TestResult, *dto.Pagination, error)
	Delete(
		ctx context.Context,
		id uint,
	) (*models.TestResult, error)
	DeleteAll(
		ctx context.Context,
		testId uuid.UUID,
	) error
}

type resultsCases struct {
	resultService ResultServiceProvider
}

func NewResultsCases(r ResultServiceProvider) *resultsCases {
	return &resultsCases{resultService: r}
}

// Работа админа c результами
func (res *resultsCases) GetListByTest(ctx context.Context, testId uuid.UUID, skip, take uint) (*dto.ResultsReponse, error) {
	results, pagination, err := res.resultService.GetByTestID(ctx, testId, skip, take)
	if err != nil {
		return nil, err
	}

	resultDtos := make([]dto.ResultReponse, 0, 100)

	for _, result := range results {
		resultDto := dto.ResultReponse{
			ID:         result.ID,
			FullName:   result.FullName,
			Email:      result.Email,
			Org:        result.Org,
			Duration:   result.Duration,
			IsOnTime:   result.IsOnTime,
			TotalScore: result.TotalScore,
		}

		resultDtos = append(resultDtos, resultDto)
	}

	return &dto.ResultsReponse{
		Data:       resultDtos,
		Pagination: *pagination,
	}, nil
}

// func (res *resultsCases) SendResultByEmail()

// func (res *resultsCases) SendAllResultsTestByEmail()

// func (res *resultsCases) MakeExcelList()

func (res *resultsCases) DeleteResult(ctx context.Context, resultId uint) (*dto.Message, error) {
	_, err := res.resultService.Delete(ctx, resultId)
	if err != nil {
		return nil, err
	}

	return &dto.Message{
		Message: "Результаты удалены",
	}, nil
}

func (res *resultsCases) DeleteResultsByTest(ctx context.Context, testId uuid.UUID) (*dto.Message, error) {
	err := res.resultService.DeleteAll(ctx, testId)
	if err != nil {
		return nil, err
	}

	return &dto.Message{
		Message: "Результаты удалены",
	}, nil
}
