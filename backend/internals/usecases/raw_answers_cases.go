package usecases

import (
	"context"
	"errors"
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
	GetByTestID(ctx context.Context, testId uuid.UUID, skip, take uint) ([]models.RawSubmission, *dto.Pagination, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.RawSubmission, error)
	СonvertFromJSON(data datatypes.JSON) (dto.SendUserAnswersRequest, error)
}

type ResultsSaveProvider interface {
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
		rawDto := dto.RawInfoResponse{
			ID:       raw.ID.String(),
			FullName: raw.FullName,
			Email:    raw.Email,
			Org:      raw.Org,
			Status:   raw.Status,
			StartAt:  raw.StartAt.Format("2006-01-02 15:04:05"),
			EndAt:    raw.EndAt.Format("2006-01-02 15:04:05"),
		}

		data = append(data, rawDto)
	}

	return &dto.RawsInfoResponse{
		Data:       data,
		Pagination: *pagination,
	}, nil
}

func (rw *rawAnswersCases) AnalyzeResults(ctx context.Context, rawId uuid.UUID) error {
	raw, err := rw.rawService.GetByID(ctx, rawId)
	if err != nil {
		// TODO: обработать
	}

	if raw.Status == "started" {
		return errors.New("Not")
	}

	var totalScore int
	var duration uint

	userAnswers, err := rw.rawService.СonvertFromJSON(raw.AnswersPayload)
	if err != nil {
		// TODO: обработать
	}

	answers := userAnswers.Answers

	for _, answer := range answers {
		userOptions := answer.OptionIDs

		question, err := rw.questionService.GetByID(ctx, answer.ID)
		if err != nil {
			// TODO: обработать
		}

		options, err := rw.optionService.GetCorrectOptions(ctx, question.ID)
		if err != nil {
			// TODO: обработать
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

	return err
}

func (rw *rawAnswersCases) AnalyzeAllResults(ctx context.Context, testId uuid.UUID)

// func (rw *rawAnswersCases) MakeReportByAnalyze(ctx context.Context, rawId uuid.UUID)

func (rw *rawAnswersCases) DeleteRawResults(ctx context.Context, rawId uuid.UUID)

func (rw *rawAnswersCases) DeleteAllRawByTest(ctx context.Context, testId uuid.UUID)
