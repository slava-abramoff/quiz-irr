package usecases

import (
	"context"
	"fmt"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"
	"sort"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type QuestionDataProvider interface {
	GetByID(ctx context.Context, id uint) (*models.Question, error)
}
type OptionDataProvider interface {
	GetCorrectOptions(
		ctx context.Context,
		questionId uint,
	) ([]models.Option, error)
}

type RawDataProvider interface {
	GetAllByTest(
		ctx context.Context,
		testId uuid.UUID,
	) ([]models.RawSubmission, error)
	GetByTestID(ctx context.Context, testId uuid.UUID, skip, take uint) ([]models.RawSubmission, *dto.Pagination, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error)
	СonvertFromJSON(data datatypes.JSON) (dto.SendUserAnswersRequest, error)
	Update(ctx context.Context, id uuid.UUID, data dto.RawUpdateRequest) (*models.RawSubmission, error)
	Delete(ctx context.Context, rawId uuid.UUID) (*models.RawSubmission, error)
	DeleteAll(ctx context.Context, testId uuid.UUID) error
}

type ResultsSaveProvider interface {
	CreateMany(ctx context.Context, results []dto.CreateResult) error
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
}

type rawAnswersCases struct {
	rawService      RawDataProvider
	questionService QuestionDataProvider
	optionService   OptionDataProvider
	resultsService  ResultsSaveProvider
}

func NewRawAnswersCases(
	r RawDataProvider,
	q QuestionDataProvider,
	o OptionDataProvider,
	res ResultsSaveProvider,
) *rawAnswersCases {
	return &rawAnswersCases{
		rawService:      r,
		questionService: q,
		optionService:   o,
		resultsService:  res,
	}
}

func (rw *rawAnswersCases) FindRawResults(ctx context.Context, testId uuid.UUID, skip, take uint) (*dto.RawsInfoResponse, error) {
	raws, pagination, err := rw.rawService.GetByTestID(ctx, testId, skip, take)
	if err != nil {
		return nil, err
	}

	data := make([]dto.RawInfoResponse, 0, 100)

	for _, raw := range raws {
		var startAt string
		if raw.StartAt != nil {
			startAt = raw.StartAt.Format("2006-01-02 15:04:05")
		}

		var endAt string
		if raw.EndAt != nil {
			endAt = raw.EndAt.Format("2006-01-02 15:04:05")
		}

		rawDto := dto.RawInfoResponse{
			ID:       raw.ID.String(),
			FullName: raw.FullName,
			Email:    raw.Email,
			Org:      raw.Org,
			Status:   raw.Status,
			StartAt:  startAt,
			EndAt:    endAt,
		}

		data = append(data, rawDto)
	}

	if pagination == nil {
		pagination = &dto.Pagination{}
	}

	return &dto.RawsInfoResponse{
		Data:       data,
		Pagination: *pagination,
	}, nil
}

func (rw *rawAnswersCases) AnalyzeResults(ctx context.Context, rawId uuid.UUID) (*dto.Message, error) {
	raw, err := rw.rawService.GetByID(ctx, rawId)
	if err != nil {
		return nil, err
	}

	if raw.Status == "started" {
		return &dto.Message{Message: "Тест находится в процессе прохождения"}, nil
	}

	if raw.Status == "handled" {
		return &dto.Message{Message: "Ответы пользователя уже обработаны"}, nil
	}

	var totalScore int
	var duration uint

	userAnswers, err := rw.rawService.СonvertFromJSON(raw.AnswersPayload)
	if err != nil {
		return nil, err
	}

	answers := userAnswers.Answers

	for _, answer := range answers {
		userOptions := answer.OptionIDs

		question, err := rw.questionService.GetByID(ctx, answer.ID)
		if err != nil {
			return nil, err
		}

		options, err := rw.optionService.GetCorrectOptions(ctx, question.ID)
		if err != nil {
			return nil, err
		}

		if question.Type == "single" {
			if options[0].ID == userOptions[0] {
				totalScore += question.Points
			}

			continue
		}

		if question.Type == "multiple" {
			var correctOptionsCount int

			if len(userOptions) != len(options) {
				continue
			}

			sort.Slice(options, func(i, j int) bool {
				return options[i].ID < options[j].ID
			})

			sort.Slice(userOptions, func(i, j int) bool {
				return userOptions[i] < userOptions[j]
			})

			for i, item := range userOptions {
				if item == options[i].ID {
					correctOptionsCount++
				}
			}

			if correctOptionsCount == len(options) {
				totalScore += question.Points
			}

			continue
		}

		if question.Type == "text" && answer.Text == options[0].Text {
			totalScore += question.Points
		}
	}

	duration = uint((raw.EndAt.Sub(*raw.StartAt).Seconds()))
	isOnTime := raw.Test.Duration > duration

	_, err = rw.resultsService.Create(
		ctx,
		raw.TestID,
		raw.FullName,
		raw.Email,
		raw.Org,
		duration,
		totalScore,
		isOnTime,
	)
	if err != nil {
		return nil, err
	}

	status := "handled"
	_, err = rw.rawService.Update(ctx, raw.ID, dto.RawUpdateRequest{
		Status: &status,
	})
	if err != nil {
		return nil, err
	}

	return &dto.Message{Message: "Ответы обработаны, подготовлены результаты"}, nil
}

func (rw *rawAnswersCases) AnalyzeAllResults(ctx context.Context, testId uuid.UUID) (*dto.Message, error) {
	raws, err := rw.rawService.GetAllByTest(ctx, testId)
	if err != nil {
		return nil, err
	}

	var count, size, countErr uint
	results := make([]dto.CreateResult, 0, len(raws))

	for _, raw := range raws {
		if raw.Status == "started" || raw.Status == "handled" {
			continue
		}
		size++

		var totalScore int
		var duration uint

		userAnswers, err := rw.rawService.СonvertFromJSON(raw.AnswersPayload)
		if err != nil {
			countErr++
			continue
		}

		answers := userAnswers.Answers

		for _, answer := range answers {
			userOptions := answer.OptionIDs

			question, err := rw.questionService.GetByID(ctx, answer.ID)
			if err != nil || question == nil {
				continue
			}
			options, err := rw.optionService.GetCorrectOptions(ctx, question.ID)
			if err != nil || len(options) == 0 {
				continue
			}

			if question.Type == "single" {
				if len(userOptions) == 0 {
					continue
				}
				if options[0].ID == userOptions[0] {
					totalScore += question.Points
				}

				continue
			}

			if question.Type == "multiple" {
				var correctOptionsCount int

				if len(userOptions) != len(options) {
					continue
				}

				sort.Slice(options, func(i, j int) bool {
					return options[i].ID < options[j].ID
				})

				sort.Slice(userOptions, func(i, j int) bool {
					return userOptions[i] < userOptions[j]
				})

				for i, item := range userOptions {
					if item == options[i].ID {
						correctOptionsCount++
					}
				}

				if correctOptionsCount == len(options) {
					totalScore += question.Points
				}

				continue
			}

			if question.Type == "text" && answer.Text == options[0].Text {
				totalScore += question.Points
			}
		}

		if raw.StartAt == nil || raw.EndAt == nil {
			countErr++
			continue
		}
		duration = uint(raw.EndAt.Sub(*raw.StartAt).Seconds())
		isOnTime := raw.Test.Duration > duration

		results = append(results, dto.CreateResult{
			TestID:     raw.TestID,
			FullName:   raw.FullName,
			Email:      raw.Email,
			Org:        raw.Org,
			Duration:   duration,
			TotalScore: totalScore,
			IsOnTime:   isOnTime,
		})

		status := "handled"
		_, _ = rw.rawService.Update(ctx, raw.ID, dto.RawUpdateRequest{
			Status: &status,
		})

		count++
	}

	err = rw.resultsService.CreateMany(ctx, results)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("Успешно обработано %d из %d", count, size)
	return &dto.Message{
		Message: msg,
	}, nil
}

// func (rw *rawAnswersCases) MakeReportByAnalyze(ctx context.Context, rawId uuid.UUID)

func (rw *rawAnswersCases) DeleteRawResults(ctx context.Context, rawId uuid.UUID) (*dto.Message, error) {
	_, err := rw.rawService.Delete(ctx, rawId)
	if err != nil {
		return nil, err
	}

	return &dto.Message{
		Message: "Ответы пользолвателя удалены",
	}, nil
}

func (rw *rawAnswersCases) DeleteAllRawByTest(ctx context.Context, testId uuid.UUID) (*dto.Message, error) {
	err := rw.rawService.DeleteAll(ctx, testId)
	if err != nil {
		return nil, err
	}

	return &dto.Message{
		Message: "Ответы пользователей удалены",
	}, nil
}
