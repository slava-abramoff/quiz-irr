package services

import (
	"context"
	"encoding/json"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type RawStorage interface {
	Create(
		ctx context.Context,
		testId uuid.UUID,
		fn, email, org string,
		start time.Time,
	) (*models.RawSubmission, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error)
	GetByTestID(ctx context.Context, testId uuid.UUID, skip, take uint) ([]models.RawSubmission, uint, error)
	SavePayload(
		ctx context.Context,
		id uuid.UUID,
		payload datatypes.JSON,
	) (*models.RawSubmission, error)
	Load(
		ctx context.Context,
		id uuid.UUID,
	) (datatypes.JSON, error)
	Update(ctx context.Context, id uuid.UUID, updated map[string]any) (*models.RawSubmission, error)
	Delete(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error)
	DeleteAll(ctx context.Context, testId uuid.UUID) error
}

type rawService struct {
	storage RawStorage
}

func NewRawService(r RawStorage) *rawService {
	return &rawService{storage: r}
}

func (rd *rawService) Create(ctx context.Context,
	testId uuid.UUID,
	fn, email, org string,
	start time.Time) (*models.RawSubmission, error) {
	return rd.storage.Create(ctx, testId, fn, email, org, start)
}

func (rd *rawService) GetByID(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error) {
	return rd.storage.GetByID(ctx, id)
}

func (rd *rawService) GetByTestID(ctx context.Context, testId uuid.UUID, skip, take uint) ([]models.RawSubmission, *dto.Pagination, error) {
	raws, count, err := rd.storage.GetByTestID(ctx, testId, skip, take)
	if err != nil {
		return raws, nil, err
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

	return raws, pagination, nil
}

func (rd *rawService) SavePayload(ctx context.Context, id uuid.UUID, data dto.SendUserAnswersRequest) (*models.RawSubmission, error) {
	payload, err := rd.СonvertToJSON(data)
	if err != nil {
		return nil, err
	}

	return rd.storage.SavePayload(ctx, id, payload)
}

func (rd *rawService) Load(ctx context.Context, id uuid.UUID) (*dto.SendUserAnswersRequest, error) {
	load, err := rd.storage.Load(ctx, id)
	if err != nil {
		return nil, err
	}

	answers, err := rd.СonvertFromJSON(load)
	if err != nil {
		return nil, err
	}

	return &answers, nil
}

func (rd *rawService) Delete(ctx context.Context, rawId uuid.UUID) (*models.RawSubmission, error) {
	return rd.storage.Delete(ctx, rawId)
}

func (rd *rawService) DeleteAll(ctx context.Context, testId uuid.UUID) error {
	return rd.storage.DeleteAll(ctx, testId)
}

func (rd *rawService) СonvertToJSON(req dto.SendUserAnswersRequest) (datatypes.JSON, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return datatypes.JSON(data), nil
}

func (rd *rawService) СonvertFromJSON(data datatypes.JSON) (dto.SendUserAnswersRequest, error) {
	var req dto.SendUserAnswersRequest

	err := json.Unmarshal(data, &req)
	if err != nil {
		return dto.SendUserAnswersRequest{}, err
	}

	return req, nil
}
