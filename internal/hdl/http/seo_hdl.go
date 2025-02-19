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

func RegisterSEORoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc(
		"/api/seo", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				middleware.Apply(h.CreateSEO, middleware.Auth(h.sso))(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)

	mux.HandleFunc(
		"/api/seo/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				h.GetSEO(w, r)
			case http.MethodPut:
				middleware.Apply(h.UpdateSEO, middleware.Auth(h.sso))(w, r)
			case http.MethodDelete:
				middleware.Apply(h.DeleteSEO, middleware.Auth(h.sso))(w, r)
			default:
				utils.ErrResponse(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			}
		},
	)
}

func (h *Handler) GetSEO(w http.ResponseWriter, r *http.Request) {
	const op = "seo.GetItemSEO.hdl"
	s, c := time.Now(), http.StatusOK
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := utils.ParseURLParams(r.URL.Path)
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.String("path", r.URL.Path),
			zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	res, err := h.ctrl.GetSEO(ctx, name, pk)
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

func (h *Handler) CreateSEO(w http.ResponseWriter, r *http.Request) {
	const op = "seo.CreateSEO.hdl"
	s, c := time.Now(), http.StatusCreated
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	req := &md.SEO{}
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

	if err := validation.ValidateSEO(req); err != nil {
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

	res, err := h.ctrl.CreateSEO(ctx, req)
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

func (h *Handler) UpdateSEO(w http.ResponseWriter, r *http.Request) {
	const op = "seo.UpdateSEO.hdl"
	s, c := time.Now(), http.StatusOK
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := utils.ParseURLParams(r.URL.Path)
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.String("path", r.URL.Path),
			zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	req := &md.SEO{}
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

	if err := validation.ValidateSEO(req); err != nil {
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

	err := h.ctrl.UpdateSEO(ctx, req)
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

func (h *Handler) DeleteSEO(w http.ResponseWriter, r *http.Request) {
	const op = "seo.DeleteSEO.hdl"
	s, c := time.Now(), http.StatusNoContent
	span, ctx := ot.StartSpanFromContext(r.Context(), op)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := utils.ParseURLParams(r.URL.Path)
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		span.SetTag("error", true)
		zap.L().Debug(
			hdl.ErrDecodeRequest.Error(),
			zap.String("op", op),
			zap.String("path", r.URL.Path),
			zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	err := h.ctrl.DeleteSEO(ctx, name, pk)
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
