package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/pkg/httpresponse"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

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
	GetTestFullData(
		ctx context.Context,
		testId uuid.UUID,
	) (*dto.TestAdminResponse, error)
	UpdateTest(
		ctx context.Context,
		data dto.UpdateTestRequest,
	) (*dto.TestAdminResponse, error)
	DeleteTest(
		ctx context.Context,
		id uuid.UUID,
	) error
	AddQuestion(ctx context.Context, testId uuid.UUID) (*dto.QuestionResponse, error)
	EditQuestion(ctx context.Context, id uint, data dto.UpdateQuestionRequest)
	DeleteQuestion(ctx context.Context, id uint) error
	AddOption(ctx context.Context, questionId uint) (*dto.OptionResponse, error)
	EditOption(ctx context.Context, id uint, data dto.UpdateOptionRequest) (*dto.OptionResponse, error)
	DeleteOption(ctx context.Context, id uint) error
}

type TestHandlers struct {
	tests TestsProvider
}

func NewTestHandlers(t TestsProvider) *TestHandlers {
	return &TestHandlers{tests: t}
}

// POST /api/tests/init
func (t *TestHandlers) NewTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var req dto.CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	// test, err := t.tests.NewTest(ctx, "DSFGsdfg", req)
	// if err != nil {
	// 	httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
	// 	return
	// }

	// httpresponse.JsonResponse(w, test, 200)
}

// GET /api/tests/:id/preview
func (t *TestHandlers) GetTestPreview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// GET /api/tests/:id/fulldata
func (t *TestHandlers) GetTestFullData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// PATCH /api/tests/:id
func (t *TestHandlers) UpdateTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// DELETE /api/tests/:id
func (t *TestHandlers) DeleteTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// POST /api/questions/test/:id
func (t *TestHandlers) AddQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// PATCH /api/questions/:id
func (t *TestHandlers) EditQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// DELETE /api/questions/:id
func (t *TestHandlers) DeleteQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// POST /api/options/question/:id
func (t *TestHandlers) AddOption(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// PATCH /api/options/:id
func (t *TestHandlers) EditOption(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}

// DELETE /api/options/:id
func (t *TestHandlers) DeleteOption(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
}
