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
	GetAllByTest(
		ctx context.Context,
		testId uuid.UUID,
	) ([]models.RawSubmission, error)
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

func (rd *rawService) GetAllByTest(
	ctx context.Context,
	testId uuid.UUID,
) ([]models.RawSubmission, error) {
	return rd.storage.GetAllByTest(ctx, testId)
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

func (rd *rawService) Update(ctx context.Context, id uuid.UUID, data dto.RawUpdateRequest) (*models.RawSubmission, error) {
	updated := make(map[string]any)

	if data.FullName != nil {
		updated["full_name"] = *data.FullName
	}
	if data.Email != nil {
		updated["email"] = *data.Email
	}
	if data.Org != nil {
		updated["org"] = *data.Org
	}
	if data.Status != nil {
		updated["status"] = *data.Status
	}

	// Accept both UI-friendly and RFC3339 timestamps.
	if data.StartAt != nil && *data.StartAt != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", *data.StartAt, time.Local); err == nil {
			tt := t
			updated["start_at"] = &tt
		} else if t, err := time.Parse(time.RFC3339, *data.StartAt); err == nil {
			tt := t
			updated["start_at"] = &tt
		} else {
			return nil, err
		}
	}
	if data.EndAt != nil && *data.EndAt != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", *data.EndAt, time.Local); err == nil {
			tt := t
			updated["end_at"] = &tt
		} else if t, err := time.Parse(time.RFC3339, *data.EndAt); err == nil {
			tt := t
			updated["end_at"] = &tt
		} else {
			return nil, err
		}
	}

	if len(updated) == 0 {
		return rd.storage.GetByID(ctx, id)
	}

	return rd.storage.Update(ctx, id, updated)
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
