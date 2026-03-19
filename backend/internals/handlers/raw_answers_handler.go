package handlers

import (
	"context"
	"log"
	"net/http"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/pkg/httpresponse"
	"strconv"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type RawAnswersProvider interface {
	FindRawResults(ctx context.Context, testId uuid.UUID, skip, take uint) (*dto.RawsInfoResponse, error)
	AnalyzeResults(ctx context.Context, rawId uuid.UUID) (*dto.Message, error)
	AnalyzeAllResults(ctx context.Context, testId uuid.UUID) (*dto.Message, error)
	DeleteRawResults(ctx context.Context, rawId uuid.UUID) (*dto.Message, error)
	DeleteAllRawByTest(ctx context.Context, testId uuid.UUID) (*dto.Message, error)
}

type RawAnswersHandlers struct {
	raw RawAnswersProvider
}

func NewRawAnswersHandlers(rw RawAnswersProvider) *RawAnswersHandlers {
	return &RawAnswersHandlers{raw: rw}
}

// GET /api/raws/test/:testId
func (rw *RawAnswersHandlers) FindRawResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	rawResults, err := rw.raw.FindRawResults(ctx, testId, uint(skip), uint(take))
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, rawResults, 200)
	log.Println("find raw results")
}

// POST /api/raws/analyze/:rawId
func (rw *RawAnswersHandlers) AnalyzeResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("start analyze results")
	ctx := r.Context()

	id := ps.ByName("rawId")
	rawId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid raw id", http.StatusBadRequest)
		return
	}

	msg, err := rw.raw.AnalyzeResults(ctx, rawId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
	log.Println("analyze results")
}

// POST /api/raws/analyze/test/:testId
func (rw *RawAnswersHandlers) AnalyzeAllResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("start full analyze results")
	ctx := r.Context()

	id := ps.ByName("testId")
	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	msg, err := rw.raw.AnalyzeAllResults(ctx, testId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
	log.Println("analyze results")
}

// DELETE /api/raws/:rawId
func (rw *RawAnswersHandlers) DeleteRawResults(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("rawId")
	rawId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid raw id", http.StatusBadRequest)
		return
	}

	msg, err := rw.raw.DeleteRawResults(ctx, rawId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
	log.Println("delete raw results")
}

// DELETE /api/raws/test/:testId
func (rw *RawAnswersHandlers) DeleteAllRawByTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("testId")
	testId, err := uuid.Parse(id)
	if err != nil {
		httpresponse.ErrorResponse(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	msg, err := rw.raw.DeleteAllRawByTest(ctx, testId)
	if err != nil {
		httpresponse.ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	httpresponse.JsonResponse(w, msg, 200)
	log.Println("delete all results")
}
