package utils

import (
	"encoding/json"
	md "github.com/JMURv/par-pro-seo/pkg/model"
	"net/http"
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
