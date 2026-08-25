package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/11DingKing/majorpath/internal/domain"
	"net/http"
)

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{Code: code, Message: message})
}
func classify(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		return 400, "invalid_request"
	case errors.Is(err, domain.ErrUnauthorized):
		return 401, "unauthorized"
	case errors.Is(err, domain.ErrForbidden):
		return 403, "forbidden"
	case errors.Is(err, domain.ErrNotFound):
		return 404, "not_found"
	case errors.Is(err, domain.ErrConflict):
		return 409, "conflict"
	default:
		return 500, "internal_error"
	}
}
