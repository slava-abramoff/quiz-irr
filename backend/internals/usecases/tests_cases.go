package usecases

import (
	"context"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/storage/models"

	"github.com/google/uuid"
)

type OptionServiceProvider interface {
	Create(ctx context.Context, questionId uint, data dto.CreateOptionRequest) (*models.Option, error)
	Update(ctx context.Context, id uint, data dto.UpdateOptionRequest) (*models.Option, error)
	GetByQuestionID(ctx context.Context, questionId uint) ([]models.Option, error)
	Delete(ctx context.Context, id uint) (*models.Option, error)
}

type QuestionServiceProvider interface {
	Create(ctx context.Context, testId uuid.UUID, data dto.CreateQuestionRequest) (*models.Question, error)
	Update(ctx context.Context, id uint, data dto.UpdateQuestionRequest) (*models.Question, error)
	GetByTestID(ctx context.Context, testId uuid.UUID) ([]models.Question, error)
	Delete(ctx context.Context, id uint) (*models.Question, error)
}

type TestServiceProvider interface {
	Create(ctx context.Context, authorId uint, data dto.CreateTestRequest) (*models.Test, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Test, error)
	FindMany(ctx context.Context, skip, take uint) ([]models.Test, *dto.Pagination, error)
	Update(ctx context.Context, id uuid.UUID, data dto.UpdateTestRequest) (*models.Test, error)
	Delete(ctx context.Context, id uuid.UUID) (*models.Test, error)
}

type testsCases struct {
	testService     TestServiceProvider
	optionService   OptionServiceProvider
	questionService QuestionServiceProvider
}

func NewTestsCases(
	t TestServiceProvider,
	o OptionServiceProvider,
	q QuestionServiceProvider,
) *testsCases {
	return &testsCases{
		testService:     t,
		questionService: q,
		optionService:   o,
	}
}

// NewTest создает новый тест в системе
func (t *testsCases) NewTest(
	ctx context.Context,
	authorId uint,
	data dto.CreateTestRequest,
) (*dto.TestAdminResponse, error) {
	test, err := t.testService.Create(ctx, authorId, data)
	if err != nil {
		return nil, err
	}

	var startAt, endAt string
	if test.StartAt != nil {
		startAt = test.StartAt.Format("2006-01-02 15:04:05")
	}
	if test.EndAt != nil {
		endAt = test.EndAt.Format("2006-01-02 15:04:05")
	}

	return &dto.TestAdminResponse{
		ID:       test.ID.String(),
		Title:    test.Title,
		Desc:     test.Desc,
		IsActive: test.IsActive,
		StartAt:  startAt,
		EndAt:    endAt,
		Author:   test.Author.FullName,
	}, nil
}

// GetTestPreview выводит информацию по тесту, не включая вопросов и ответов
func (t *testsCases) GetTestPreview(
	ctx context.Context,
	testId uuid.UUID,
) (*dto.TestAdminResponse, error) {
	test, err := t.testService.GetByID(ctx, testId)
	if err != nil {
		return nil, err
	}

	var startAt, endAt string
	if test.StartAt != nil {
		startAt = test.StartAt.Format("2006-01-02 15:04:05")
	}
	if test.EndAt != nil {
		endAt = test.EndAt.Format("2006-01-02 15:04:05")
	}

	return &dto.TestAdminResponse{
		ID:       test.ID.String(),
		Title:    test.Title,
		Desc:     test.Desc,
		IsActive: test.IsActive,
		StartAt:  startAt,
		EndAt:    endAt,
		Duration: test.Duration,
		Author:   test.Author.FullName,
	}, nil
}

func (t *testsCases) FindManyTests(ctx context.Context, skip, take uint) (*dto.GetManyTestsResponse, error) {
	tests, pagination, err := t.testService.FindMany(ctx, skip, take)
	if err != nil {
		return nil, err
	}

	testDtos := make([]dto.TestAdminResponse, 0, take)

	for _, test := range tests {
		var startAt, endAt string
		if test.StartAt != nil {
			startAt = test.StartAt.Format("2006-01-02 15:04:05")
		}
		if test.EndAt != nil {
			endAt = test.EndAt.Format("2006-01-02 15:04:05")
		}

		testDto := dto.TestAdminResponse{
			ID:       test.ID.String(),
			Title:    test.Title,
			Desc:     test.Desc,
			IsActive: test.IsActive,
			StartAt:  startAt,
			EndAt:    endAt,
			Author:   test.Author.FullName,
			Duration: test.Duration,
		}

		testDtos = append(testDtos, testDto)
	}

	return &dto.GetManyTestsResponse{
		Pagination: *pagination,
		Tests:      testDtos,
	}, nil
}

// GetTestFullData выводит информацию по тесту, включая вопросов и ответов
func (t *testsCases) GetTestFullData(
	ctx context.Context,
	testId uuid.UUID,
) (*dto.TestAdminResponse, error) {
	var testData dto.TestAdminResponse

	test, err := t.testService.GetByID(ctx, testId)
	if err != nil {
		return nil, err
	}

	var startAt, endAt string
	if test.StartAt != nil {
		startAt = test.StartAt.Format("2006-01-02 15:04:05")
	}
	if test.EndAt != nil {
		endAt = test.EndAt.Format("2006-01-02 15:04:05")
	}

	testData = dto.TestAdminResponse{
		ID:       test.ID.String(),
		Title:    test.Title,
		Desc:     test.Desc,
		IsActive: test.IsActive,
		StartAt:  startAt,
		EndAt:    endAt,
		Author:   test.Author.FullName,
	}

	questions, err := t.questionService.GetByTestID(ctx, test.ID)
	if err != nil {
		return nil, err
	}

	questionDtos := make([]dto.QuestionResponse, 0, 75)

	for _, q := range questions {
		optionDtos := make([]dto.OptionResponse, 0, 10)

		questionDto := dto.QuestionResponse{
			ID:     q.ID,
			Text:   q.Text,
			Type:   q.Type,
			Points: q.Points,
		}

		options, err := t.optionService.GetByQuestionID(ctx, q.ID)
		if err != nil {
			return nil, err
		}

		for _, o := range options {
			optionDto := dto.OptionResponse{
				ID:        o.ID,
				Text:      o.Text,
				IsCorrect: o.IsCorrect,
			}

			optionDtos = append(optionDtos, optionDto)
		}

		questionDto.Options = optionDtos
		questionDtos = append(questionDtos, questionDto)
	}

	testData.Questions = questionDtos
	return &testData, nil
}

// UpdateTest обновляет некоторые поля теста
func (t *testsCases) UpdateTest(
	ctx context.Context,
	id uuid.UUID,
	data dto.UpdateTestRequest,
) (*dto.TestAdminResponse, error) {
	test, err := t.testService.Update(ctx, id, data)
	if err != nil {
		return nil, nil
	}

	return &dto.TestAdminResponse{
		ID:       test.ID.String(),
		Title:    test.Title,
		Desc:     test.Desc,
		IsActive: test.IsActive,
		Duration: test.Duration,
		StartAt:  test.StartAt.Format("2006-01-02 15:04:05"),
		EndAt:    test.EndAt.Format("2006-01-02 15:04:05"),
		Author:   test.Author.FullName,
	}, nil
}

// DeleteTest удаляет конкретный тест
func (t *testsCases) DeleteTest(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := t.testService.Delete(ctx, id)
	return err
}

// AddQuestion добавляет новый вопрос для теста
func (t *testsCases) AddQuestion(ctx context.Context, testId uuid.UUID) (*dto.QuestionResponse, error) {
	question, err := t.questionService.Create(ctx, testId, dto.CreateQuestionRequest{
		Text:   "Сотрите текст и напишите свой вопрос",
		Type:   "multiple",
		Points: 1,
	})
	if err != nil {
		return nil, err
	}

	return &dto.QuestionResponse{
		ID:     question.ID,
		Text:   question.Text,
		Type:   question.Type,
		Points: question.Points,
	}, nil
}

// EditQuestion позволяет редактировать текстовое поле и тип вопроса
func (t *testsCases) EditQuestion(ctx context.Context, id uint, data dto.UpdateQuestionRequest) (*dto.QuestionResponse, error) {
	question, err := t.questionService.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}

	options, err := t.optionService.GetByQuestionID(ctx, question.ID)
	if err != nil {
		return nil, err
	}

	optionDtos := make([]dto.OptionResponse, 0, 10)

	for _, o := range options {
		optionDto := dto.OptionResponse{
			ID:        o.ID,
			Text:      o.Text,
			IsCorrect: o.IsCorrect,
		}

		optionDtos = append(optionDtos, optionDto)
	}

	return &dto.QuestionResponse{
		ID:      question.ID,
		Text:    question.Text,
		Type:    question.Type,
		Points:  question.Points,
		Options: optionDtos,
	}, nil
}

// DeleteQuestion удаляет конкретный вопрос
func (t *testsCases) DeleteQuestion(ctx context.Context, id uint) error {
	_, err := t.questionService.Delete(ctx, id)
	return err
}

// AddOption добавляет новый ответ для вопроса
func (t *testsCases) AddOption(ctx context.Context, questionId uint) (*dto.OptionResponse, error) {
	option, err := t.optionService.Create(ctx, questionId, dto.CreateOptionRequest{
		Text:      "Новый ответ",
		IsCorrect: false,
	})
	if err != nil {
		return nil, err
	}

	return &dto.OptionResponse{
		ID:        option.ID,
		Text:      option.Text,
		IsCorrect: option.IsCorrect,
	}, nil
}

// EditOption позволяет редактировать текстовое поле и корректность ответа
func (t *testsCases) EditOption(ctx context.Context, id uint, data dto.UpdateOptionRequest) (*dto.OptionResponse, error) {
	option, err := t.optionService.Update(ctx, id, data)
	if err != nil {
		return nil, err
	}

	return &dto.OptionResponse{
		ID:        option.ID,
		Text:      option.Text,
		IsCorrect: option.IsCorrect,
	}, nil
}

// DeleteOption удаляет конкретный ответ
func (t *testsCases) DeleteOption(ctx context.Context, id uint) error {
	_, err := t.optionService.Delete(ctx, id)
	return err
}
