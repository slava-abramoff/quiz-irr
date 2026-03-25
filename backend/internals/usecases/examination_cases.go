package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"quiz-irr/pkg/apperrors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	Create(ctx context.Context, testId uuid.UUID, f, e, o string, b uint, start time.Time) (*models.RawSubmission, error)
}

type examCases struct {
	testService     TestInfoProvider
	optionService   OptionInfoProvider
	questionService QuestionInfoProvider
	rawDataService  RawDataServiceProvider
	cache           *redis.Client
}

func NewExamCases(t TestInfoProvider, r RawDataServiceProvider, o OptionInfoProvider, q QuestionInfoProvider, cache *redis.Client) *examCases {
	return &examCases{
		testService:     t,
		rawDataService:  r,
		optionService:   o,
		questionService: q,
		cache:           cache,
	}
}

// Получить информацию о тесте
func (e *examCases) GetTestInfo(ctx context.Context, id uuid.UUID) (*dto.TestCustomerResponse, error) {
	test, err := e.testService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !test.IsActive {
		return nil, apperrors.NotFound("Test not found")
	}

	return &dto.TestCustomerResponse{
		Title:    test.Title,
		Desc:     test.Desc,
		Duration: test.Duration,
		StartAt:  test.StartAt.Format(time.RFC3339),
		EndAt:    test.EndAt.Format(time.RFC3339),
	}, nil
}

func (e *examCases) StartTest(ctx context.Context, id uuid.UUID, data dto.StartExamRequest) (*dto.StartExamResponse, error) {
	test, err := e.testService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if test.StartAt != nil && test.EndAt != nil {
		if !e.isTimeToStart(*test.StartAt, *test.EndAt) {
			return nil, apperrors.Conflict("Test is not available at this time")
		}
	}

	cacheKey := fmt.Sprintf("exam:test:%s:questions:v1", id.String())
	var questionDtos []dto.ExamQuestion
	cacheHit := false

	if e.cache != nil {
		cached, err := e.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cached), &questionDtos); err == nil {
				cacheHit = true
			}
		} else if !errors.Is(err, redis.Nil) {
			// cache error: do not fail the request, just fallback to DB
		}
	}

	if cacheHit {
		log.Println("Cache hit")
	} else {
		if e.cache == nil {
			log.Println("Cache disabled (redis not connected)")
		} else {
			log.Println("Cache miss")
		}
		questions, err := e.questionService.GetByTestID(ctx, id)
		if err != nil {
			return nil, err
		}

		questionDtos = make([]dto.ExamQuestion, 0, len(questions))

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

		if e.cache != nil {
			if b, err := json.Marshal(questionDtos); err == nil {
				if err := e.cache.Set(ctx, cacheKey, b, 30*time.Minute).Err(); err == nil {
					log.Println("Cached")
				}
			}
		}
	}

	startAt := time.Now()

	raw, err := e.rawDataService.Create(
		ctx,
		id,
		data.FullName,
		data.Email,
		data.Org,
		data.BirthYear,
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
