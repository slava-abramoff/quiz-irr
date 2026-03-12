package services

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"time"

	"github.com/google/uuid"
)

type TestStorageProvider interface {
	Create(
		ctx context.Context,
		authorId uint,
		title string,
		desc string,
		isActive bool,
		startAt time.Time,
		endAt time.Time,
	) (*models.Test, error)
	FindMany(ctx context.Context, skip, take uint) ([]models.Test, uint, error)
	GetByAuthor(ctx context.Context, authorId uint) ([]models.Test, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Test, error)
	Update(
		ctx context.Context,
		id uuid.UUID,
		updated map[string]any,
	) (*models.Test, error)
	Delete(ctx context.Context, id uuid.UUID) (*models.Test, error)
}

type testService struct {
	storage TestStorageProvider
}

func NewTestService(s TestStorageProvider) *testService {
	return &testService{storage: s}
}

func (t *testService) Create(ctx context.Context, authorId uint, data dto.CreateTestRequest) (*models.Test, error) {
	return t.storage.Create(
		ctx,
		authorId,
		data.Title,
		data.Desc,
		false,
		time.Now(),
		time.Now(),
	)
}

func (t *testService) GetByID(ctx context.Context, id uuid.UUID) (*models.Test, error) {
	return t.storage.GetByID(ctx, id)
}

func (t *testService) FindMany(ctx context.Context, skip, take uint) ([]models.Test, *dto.Pagination, error) {
	tests, count, err := t.storage.FindMany(ctx, skip, take)
	if err != nil {
		return tests, nil, err
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

	return tests, pagination, nil
}

func (t *testService) Update(ctx context.Context, id uuid.UUID, data dto.UpdateTestRequest) (*models.Test, error) {
	updated := make(map[string]any)

	if data.Title != nil {
		updated["title"] = *data.Title
	}

	if data.Desc != nil {
		updated["desc"] = *data.Desc
	}

	if data.IsActive != nil {
		updated["is_active"] = *data.IsActive
	}

	if data.Duration != nil {
		updated["duration"] = *data.Duration
	}

	if data.StartAt != nil {
		updated["start_at"] = *data.StartAt
	}

	if data.EndAt != nil {
		updated["title"] = *data.Title
	}

	return t.storage.Update(ctx, id, updated)
}

func (t *testService) Delete(ctx context.Context, id uuid.UUID) (*models.Test, error) {
	return t.storage.Delete(ctx, id)
}
