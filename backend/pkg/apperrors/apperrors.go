package apperrors

import (
	"errors"
	"net/http"

	"gorm.io/gorm"
)

// Sentinel errors for errors.Is() checks.
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation failed")
	ErrBadRequest   = errors.New("bad request")
)

// withCode carries HTTP code and user-facing message.
type withCode struct {
	code int
	msg  string
	err  error
}

func (e *withCode) Error() string {
	if e.msg != "" {
		return e.msg
	}
	if e.err != nil {
		return e.err.Error()
	}
	return "error"
}

func (e *withCode) Unwrap() error { return e.err }

// NotFound returns error that maps to 404.
func NotFound(msg string) error {
	if msg == "" {
		msg = "Not found"
	}
	return &withCode{code: http.StatusNotFound, msg: msg, err: ErrNotFound}
}

// Unauthorized returns error that maps to 401.
func Unauthorized(msg string) error {
	if msg == "" {
		msg = "Unauthorized"
	}
	return &withCode{code: http.StatusUnauthorized, msg: msg, err: ErrUnauthorized}
}

// Forbidden returns error that maps to 403.
func Forbidden(msg string) error {
	if msg == "" {
		msg = "Forbidden"
	}
	return &withCode{code: http.StatusForbidden, msg: msg, err: ErrForbidden}
}

// Conflict returns error that maps to 409.
func Conflict(msg string) error {
	if msg == "" {
		msg = "Conflict"
	}
	return &withCode{code: http.StatusConflict, msg: msg, err: ErrConflict}
}

// Validation returns error that maps to 422.
func Validation(msg string) error {
	if msg == "" {
		msg = "Validation failed"
	}
	return &withCode{code: http.StatusUnprocessableEntity, msg: msg, err: ErrValidation}
}

// BadRequest returns error that maps to 400.
func BadRequest(msg string) error {
	if msg == "" {
		msg = "Bad request"
	}
	return &withCode{code: http.StatusBadRequest, msg: msg, err: ErrBadRequest}
}

// ToHTTP maps an error to HTTP status code and safe user-facing message.
// Unknown errors get 500 and generic "Internal server error".
func ToHTTP(err error) (code int, message string) {
	if err == nil {
		return 0, ""
	}

	var w *withCode
	if errors.As(err, &w) {
		return w.code, w.msg
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http.StatusNotFound, "Not found"
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound, "Not found"
	}
	if errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized, "Unauthorized"
	}
	if errors.Is(err, ErrForbidden) {
		return http.StatusForbidden, "Forbidden"
	}
	if errors.Is(err, ErrConflict) {
		return http.StatusConflict, "Conflict"
	}
	if errors.Is(err, ErrValidation) {
		return http.StatusUnprocessableEntity, "Validation failed"
	}
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest, "Bad request"
	}

	return http.StatusInternalServerError, "Internal server error"
}
