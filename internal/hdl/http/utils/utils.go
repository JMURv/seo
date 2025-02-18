package utils

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func StatusResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

func SuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func ErrResponse(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(
		&ErrorResponse{
			Error: err.Error(),
		},
	)
}

func ParseURLParams(path string) (string, string) {
	parts := strings.Split(
		strings.TrimPrefix(path, "/api/seo/"), "/",
	)

	if len(parts) != 2 {
		zap.L().Debug(
			"failed to decode request, incorrect path format",
			zap.String("path", path),
		)
		return "", ""
	}

	return parts[0], parts[1]
}

func ParsePageParams(path string) string {
	parts := strings.Split(
		strings.TrimPrefix(path, "/api/page/"), "/",
	)

	return parts[0]
}
