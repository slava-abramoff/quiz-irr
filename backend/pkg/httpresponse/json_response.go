package httpresponse

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}
