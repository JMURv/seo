package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/JMURv/seo-svc/internal/controller/sso"
	"github.com/JMURv/seo-svc/internal/handler"
	mid "github.com/JMURv/seo-svc/internal/handler/http/middleware"
	utils "github.com/JMURv/seo-svc/pkg/utils/http"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	srv  *http.Server
	ctrl handler.SEOCtrl
	sso  sso.SSOSvc
}

func New(ctrl handler.SEOCtrl, sso sso.SSOSvc) *Handler {
	return &Handler{
		ctrl: ctrl,
		sso:  sso,
	}
}

func (h *Handler) Start(port int) {
	mux := http.NewServeMux()
	RegisterSEORoutes(mux, h)
	RegisterPageRoutes(mux, h)
	mux.HandleFunc(
		"/api/health-check", func(w http.ResponseWriter, r *http.Request) {
			utils.SuccessResponse(w, http.StatusOK, "OK")
		},
	)

	hdl := mid.RecoverPanic(mux)
	hdl = mid.TracingMiddleware(mux)
	h.srv = &http.Server{
		Handler:      hdl,
		Addr:         fmt.Sprintf(":%v", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	err := h.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		zap.L().Debug("Server error", zap.Error(err))
	}
}

func (h *Handler) Close() error {
	if err := h.srv.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrResponse(w, http.StatusUnauthorized, errors.New("authorization header is missing"))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == authHeader {
				utils.ErrResponse(w, http.StatusUnauthorized, errors.New("invalid token format"))
				return
			}

			token, err := h.sso.ParseClaims(r.Context(), tokenStr)
			if err != nil {
				utils.ErrResponse(w, http.StatusUnauthorized, err)
				return
			}

			ctx := context.WithValue(r.Context(), "uid", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
