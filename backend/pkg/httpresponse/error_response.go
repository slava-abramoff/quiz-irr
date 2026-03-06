package httpresponse

import "net/http"

func ErrorResponse(w http.ResponseWriter, message string, code int) error {
	err := JsonResponse(w, map[string]any{
		"status":  "error",
		"message": message,
	}, code)
	if err != nil {
		return err
	}

	return nil
}
