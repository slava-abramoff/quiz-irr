package services

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"

	"github.com/google/uuid"
)

type ResultStorageProvider interface {
	Create(
		ctx context.Context,
		testId uuid.UUID,
		fullName string,
		email string,
		org string,
		duration uint,
		totalScore int,
		isOnTime bool,
	) (*models.TestResult, error)
	GetByID(
		ctx context.Context,
		id uint,
	) (*models.TestResult, error)
	GetAllByTestID(
		ctx context.Context,
		testId uuid.UUID,
	) ([]models.TestResult, error)
	GetByTestID(
		ctx context.Context,
		testId uuid.UUID,
		skip, take uint,
	) ([]models.TestResult, uint, error)
	Update(
		ctx context.Context,
		id uint,
		updates map[string]any,
	) (*models.TestResult, error)
	Delete(
		ctx context.Context,
		id uint,
	) (*models.TestResult, error)
	DeleteByTestID(
		ctx context.Context,
		testId uuid.UUID,
	) error
}

type resultsService struct {
	storage ResultStorageProvider
}

func NewResultsService(r ResultStorageProvider) *resultsService {
	return &resultsService{storage: r}
}

func (r *resultsService) Create(
	ctx context.Context,
	testId uuid.UUID,
	fullName string,
	email string,
	org string,
	duration uint,
	totalScore int,
	isOnTime bool,
) (*models.TestResult, error) {
	return r.storage.Create(
		ctx,
		testId,
		fullName,
		email,
		org,
		duration,
		totalScore,
		isOnTime,
	)
}

func (r *resultsService) GetAll(
	ctx context.Context,
	testId uuid.UUID,
) ([]models.TestResult, error) {
	return r.storage.GetAllByTestID(ctx, testId)
}

func (r *resultsService) GetByTestID(
	ctx context.Context,
	testId uuid.UUID,
	skip, take uint,
) ([]models.TestResult, *dto.Pagination, error) {
	results, count, err := r.storage.GetByTestID(ctx, testId, skip, take)
	if err != nil {
		return results, nil, err
	}

	currentPage := skip/take + 1
	totalPages := (count + take - 1) / take

	pagination := &dto.Pagination{
		CurrentPage:     currentPage,
		TotalPages:      totalPages,
		TotalItems:      count,
		ItemsPerPage:    take,
		HasNextPage:     currentPage < totalPages,
		HasPreviousPage: currentPage > 1,
	}

	return results, pagination, nil
}

func (r *resultsService) Delete(
	ctx context.Context,
	id uint,
) (*models.TestResult, error) {
	return r.storage.Delete(ctx, id)
}

func (r *resultsService) DeleteAll(
	ctx context.Context,
	testId uuid.UUID,
) error {
	return r.storage.DeleteByTestID(ctx, testId)
}
