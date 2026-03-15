package services

import (
	"context"
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
	GetByTestID(
		ctx context.Context,
		testId uuid.UUID,
	) ([]models.TestResult, error)
	Update(
		ctx context.Context,
		id uint,
		updates map[string]any,
	) (*models.TestResult, error)
	Delete(
		ctx context.Context,
		id uint,
	) (*models.TestResult, error)
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
