package handlers

import (
	"context"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/pkg/httpresponse"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type ResultsCasesProvider interface {
	GetListByTest(ctx context.Context, testId uuid.UUID, skip, take uint) (*dto.ResultsReponse, error)
	DeleteResult(ctx context.Context, resultId uint) (*dto.Message, error)
	DeleteResultsByTest(ctx context.Context, testId uuid.UUID) (*dto.Message, error)
}

type ResultsHandlers struct {
	results ResultsCasesProvider
}

func NewResultsHandlers(r ResultsCasesProvider) *ResultsHandlers {
	return &ResultsHandlers{results: r}
}

// GET /api/results/test/:testId
func (res *ResultsHandlers) GetListByTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("testId")
	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

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

	results, err := res.results.GetListByTest(ctx, testId, uint(skip), uint(take))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, results, 200)
}

// DELETE /api/results/:resultId
func (res *ResultsHandlers) DeleteResult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("resultId")
	resultId, err := strconv.Atoi(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid result id", http.StatusBadRequest)
		return
	}

	msg, err := res.results.DeleteResult(ctx, uint(resultId))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
}

// DELETE /api/results/test/:testId
func (res *ResultsHandlers) DeleteResultsByTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("testId")
	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	msg, err := res.results.DeleteResultsByTest(ctx, testId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
}
