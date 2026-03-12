package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/pkg/httpresponse"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type AuthInfoProvider interface {
	GetAuthInfo(r *http.Request) (uint, bool, error)
}

type TestsProvider interface {
	NewTest(
		ctx context.Context,
		authorId uint,
		data dto.CreateTestRequest,
	) (*dto.TestAdminResponse, error)
	GetTestPreview(
		ctx context.Context,
		testId uuid.UUID,
	) (*dto.TestAdminResponse, error)
	FindManyTests(ctx context.Context, skip, take uint) (*dto.GetManyTestsResponse, error)
	GetTestFullData(
		ctx context.Context,
		testId uuid.UUID,
	) (*dto.TestAdminResponse, error)
	UpdateTest(
		ctx context.Context,
		id uuid.UUID,
		data dto.UpdateTestRequest,
	) (*dto.TestAdminResponse, error)
	DeleteTest(
		ctx context.Context,
		id uuid.UUID,
	) error
	AddQuestion(ctx context.Context, testId uuid.UUID) (*dto.QuestionResponse, error)
	EditQuestion(ctx context.Context, id uint, data dto.UpdateQuestionRequest) (*dto.QuestionResponse, error)
	DeleteQuestion(ctx context.Context, id uint) error
	AddOption(ctx context.Context, questionId uint) (*dto.OptionResponse, error)
	EditOption(ctx context.Context, id uint, data dto.UpdateOptionRequest) (*dto.OptionResponse, error)
	DeleteOption(ctx context.Context, id uint) error
}

type TestHandlers struct {
	tests TestsProvider
	users AuthInfoProvider
}

func NewTestHandlers(t TestsProvider, u AuthInfoProvider) *TestHandlers {
	return &TestHandlers{tests: t, users: u}
}

// POST /api/tests/init
func (t *TestHandlers) NewTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var req dto.CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	id, _, err := t.users.GetAuthInfo(r)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid token", http.StatusForbidden)
		return
	}

	test, err := t.tests.NewTest(ctx, id, req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, test, 200)
}

// GET /api/tests/many
func (t *TestHandlers) FindManyTests(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	q := r.URL.Query()

	skipStr := q.Get("skip")
	takeStr := q.Get("take")

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid skip value", http.StatusBadRequest)
		return
	}
	take, err := strconv.Atoi(takeStr)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid take value", http.StatusBadRequest)
		return
	}

	if take == 0 {
		take += 1
	}

	tests, err := t.tests.FindManyTests(ctx, uint(skip), uint(take))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, tests, 200)
}

// GET /api/tests/:id/:mode
func (t *TestHandlers) GetTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	mode := ps.ByName("mode")
	id := ps.ByName("id")

	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	var test *dto.TestAdminResponse

	if mode == "preview" {
		test, err = t.tests.GetTestPreview(ctx, testId)
	} else {
		test, err = t.tests.GetTestFullData(ctx, testId)
	}

	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, test, 200)
}

// PATCH /api/tests/:id
func (t *TestHandlers) UpdateTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("id")

	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	test, err := t.tests.UpdateTest(ctx, testId, req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, test, 200)
}

// DELETE /api/tests/:id
func (t *TestHandlers) DeleteTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("id")

	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	err = t.tests.DeleteTest(ctx, testId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, "Deleted", 200)
}

// POST /api/questions/test/:id
func (t *TestHandlers) AddQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("id")

	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	question, err := t.tests.AddQuestion(ctx, testId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, question, 200)
}

// PATCH /api/questions/:id
func (t *TestHandlers) EditQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	questionId, err := strconv.Atoi(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid question id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	question, err := t.tests.EditQuestion(ctx, uint(questionId), req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, question, 200)
}

// DELETE /api/questions/:id
func (t *TestHandlers) DeleteQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	questionId, err := strconv.Atoi(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid question id", http.StatusBadRequest)
		return
	}

	err = t.tests.DeleteQuestion(ctx, uint(questionId))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, "Deleted", 200)
}

// POST /api/options/question/:id
func (t *TestHandlers) AddOption(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	questionId, err := strconv.Atoi(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid question id", http.StatusBadRequest)
		return
	}

	option, err := t.tests.AddOption(ctx, uint(questionId))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, option, 200)
}

// PATCH /api/options/:id
func (t *TestHandlers) EditOption(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	optionId, err := strconv.Atoi(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid option id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	option, err := t.tests.EditOption(ctx, uint(optionId), req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, option, 200)
}

// DELETE /api/options/:id
func (t *TestHandlers) DeleteOption(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	optionId, err := strconv.Atoi(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid option id", http.StatusBadRequest)
		return
	}

	err = t.tests.DeleteOption(ctx, uint(optionId))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, "Deleted", 200)
}
