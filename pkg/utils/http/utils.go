package utils

import (
	"encoding/json"
	md "github.com/JMURv/seo-svc/pkg/model"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Response struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginatedData struct {
	Data        []*md.SEO `json:"data"`
	Count       int64     `json:"count"`
	TotalPages  int       `json:"total_pages"`
	CurrentPage int       `json:"current_page"`
	HasNextPage bool      `json:"has_next_page"`
}

func SuccessPaginatedResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func SuccessResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&Response{
		Data: data,
	})
}

func ErrResponse(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorResponse{
		Error: err.Error(),
	})
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
