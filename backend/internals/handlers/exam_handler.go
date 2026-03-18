package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/pkg/httpresponse"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ExamProvider interface {
	GetTestInfo(ctx context.Context, id uuid.UUID) (*dto.TestCustomerResponse, error)
	StartTest(ctx context.Context, id uuid.UUID, data dto.StartExamRequest) (*dto.StartExamResponse, error)
	SaveAnswers(ctx context.Context, rawId uuid.UUID, data dto.SendUserAnswersRequest) (*dto.SendUserAnswersResponse, error)
}

type ExamHandlers struct {
	exam ExamProvider
}

func NewExamHandlers(e ExamProvider) *ExamHandlers {
	return &ExamHandlers{exam: e}
}

// /api/exam/info/:testId
func (e *ExamHandlers) GetTestInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("testId")

	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	info, err := e.exam.GetTestInfo(ctx, testId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, info, 200)
	log.Println("get test info")
}

// /api/exam/start/:testId
func (e *ExamHandlers) StartTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("testId")

	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	var req dto.StartExamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || len([]byte(req.Email)) > 30 {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.FullName == "" || len([]byte(req.FullName)) > 75 {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.Org == "" || len([]byte(req.Org)) > 75 {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	testBody, err := e.exam.StartTest(ctx, testId, req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, testBody, 200)
	log.Println("start test")
}

// /api/exam/save/:rawId
func (e *ExamHandlers) SaveAnswers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("rawId")

	rawId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid raw id", http.StatusBadRequest)
		return
	}

	var req dto.SendUserAnswersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	msg, err := e.exam.SaveAnswers(ctx, rawId, req)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
	log.Println("save answers")
}
