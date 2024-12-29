package middleware

import (
	"errors"
	ctrl "github.com/JMURv/seo-svc/internal/controller"
	utils "github.com/JMURv/seo-svc/pkg/utils/http"
	"go.uber.org/zap"
	"net/http"
)

var ErrMethodNotAllowed = errors.New("method not allowed")

func MethodNotAllowed(methods ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, method := range methods {
		allowed[method] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if _, ok := allowed[r.Method]; !ok {
					utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error("panic", zap.Any("err", err))
					utils.ErrResponse(w, http.StatusInternalServerError, ctrl.ErrInternalError)
				}
			}()
			next.ServeHTTP(w, r)
		},
	)
}
