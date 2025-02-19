package http

import (
	"encoding/json"
	"errors"
	"github.com/JMURv/seo/internal/ctrl"
	"github.com/JMURv/seo/internal/hdl"
	"github.com/JMURv/seo/internal/hdl/http/middleware"
	"github.com/JMURv/seo/internal/hdl/http/utils"
	"github.com/JMURv/seo/internal/hdl/validation"
	md "github.com/JMURv/seo/internal/models"
	metrics "github.com/JMURv/seo/internal/observability/metrics/prometheus"
	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegisterPageRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc(
		"/api/page", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.ListPages(w, r)
			case http.MethodPost:
				middleware.Apply(h.CreatePage, middleware.Auth(h.sso))(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)

	mux.HandleFunc(
		"/api/page/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.GetPage(w, r)
			case http.MethodPut:
				middleware.Apply(h.UpdatePage, middleware.Auth(h.sso))(w, r)
			case http.MethodDelete:
				middleware.Apply(h.DeletePage, middleware.Auth(h.sso))(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)
}

func (h *Handler) ListPages(w http.ResponseWriter, r *http.Request) {
	const op = "pages.ListPages.hdl"
	s, c := time.Now(), http.StatusOK
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	res, err := h.ctrl.ListPages(ctx)
	if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) GetPage(w http.ResponseWriter, r *http.Request) {
	const op = "pages.GetPage.hdl"
	s, c := time.Now(), http.StatusOK
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	slug := utils.ParsePageParams(r.URL.Path)
	if slug == "" {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.String("slug", slug),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	res, err := h.ctrl.GetPage(ctx, slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) CreatePage(w http.ResponseWriter, r *http.Request) {
	const op = "pages.CreatePage.hdl"
	s, c := time.Now(), http.StatusCreated
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	req := &md.Page{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidatePage(req); err != nil {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			"failed to validate",
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		utils.ErrResponse(w, c, err)
		return
	}

	res, err := h.ctrl.CreatePage(ctx, req)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		c = http.StatusConflict
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, res)
}

func (h *Handler) UpdatePage(w http.ResponseWriter, r *http.Request) {
	const op = "pages.UpdatePage.hdl"
	s, c := time.Now(), http.StatusOK
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	slug := utils.ParsePageParams(r.URL.Path)
	if slug == "" {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.String("path", r.URL.Path),
			zap.String("slug", slug),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	req := &md.Page{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidatePage(req); err != nil {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			"failed to validate",
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		utils.ErrResponse(w, c, err)
		return
	}

	err := h.ctrl.UpdatePage(ctx, slug, req)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.StatusResponse(w, c)
}

func (h *Handler) DeletePage(w http.ResponseWriter, r *http.Request) {
	const op = "pages.DeletePage.hdl"
	s, c := time.Now(), http.StatusNoContent
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	slug := utils.ParsePageParams(r.URL.Path)
	if slug == "" {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.String("path", r.URL.Path),
			zap.String("slug", slug),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	err := h.ctrl.DeletePage(ctx, slug)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.StatusResponse(w, c)
}
