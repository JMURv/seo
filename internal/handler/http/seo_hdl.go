package http

import (
	"errors"
	ctrl "github.com/JMURv/par-pro-seo/internal/controller"
	hdl "github.com/JMURv/par-pro-seo/internal/handler"
	metrics "github.com/JMURv/par-pro-seo/internal/metrics/prometheus"
	"github.com/JMURv/par-pro-seo/internal/validation"
	"github.com/JMURv/par-pro-seo/pkg/model"
	utils "github.com/JMURv/par-pro-seo/pkg/utils/http"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RegisterSEORoutes(r *mux.Router, h *Handler) {
	r.HandleFunc("/api/seo/{name}/{pk}", h.GetSEO).Methods(http.MethodGet)
	r.HandleFunc("/api/seo/{name}/{pk}", middlewareFunc(h.CreateSEO, h.authMiddleware)).Methods(http.MethodPost)
	r.HandleFunc("/api/seo/{name}/{pk}", middlewareFunc(h.UpdateSEO, h.authMiddleware)).Methods(http.MethodPut)
	r.HandleFunc("/api/seo/{name}/{pk}", middlewareFunc(h.DeleteSEO, h.authMiddleware)).Methods(http.MethodDelete)
}

func (h *Handler) GetSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "seo.GetItemSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := mux.Vars(r)["name"], mux.Vars(r)["pk"]
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request",
			zap.String("op", op), zap.String("objName", name), zap.String("objPK", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	res, err := h.ctrl.GetSEO(r.Context(), name, pk)
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
	s, c := time.Now(), http.StatusOK
	const op = "seo.CreateSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := mux.Vars(r)["name"], mux.Vars(r)["pk"]
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request",
			zap.String("op", op), zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	req := &model.SEO{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidateSEO(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	err := h.ctrl.UpdateSEO(r.Context(), name, pk, req)
	if err != nil && errors.Is(err, ctrl.ErrAlreadyExists) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) UpdateSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "seo.UpdateSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := mux.Vars(r)["name"], mux.Vars(r)["pk"]
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request",
			zap.String("op", op), zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	req := &model.SEO{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to decode request", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	if err := validation.ValidateSEO(req); err != nil {
		c = http.StatusBadRequest
		zap.L().Debug("failed to validate obj", zap.String("op", op), zap.Error(err))
		utils.ErrResponse(w, c, err)
		return
	}

	err := h.ctrl.UpdateSEO(r.Context(), name, pk, req)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}

func (h *Handler) DeleteSEO(w http.ResponseWriter, r *http.Request) {
	s, c := time.Now(), http.StatusOK
	const op = "seo.DeleteSEO.handler"
	defer func() {
		metrics.ObserveRequest(time.Since(s), c, op)
	}()

	name, pk := mux.Vars(r)["name"], mux.Vars(r)["pk"]
	if name == "" || pk == "" {
		c = http.StatusBadRequest
		zap.L().Debug(
			"failed to decode request",
			zap.String("op", op), zap.String("name", name), zap.String("pk", pk),
		)
		utils.ErrResponse(w, c, hdl.ErrDecodeRequest)
		return
	}

	err := h.ctrl.DeleteSEO(r.Context(), name, pk)
	if err != nil && errors.Is(err, ctrl.ErrNotFound) {
		c = http.StatusNotFound
		utils.ErrResponse(w, c, err)
		return
	} else if err != nil {
		c = http.StatusInternalServerError
		utils.ErrResponse(w, c, hdl.ErrInternal)
		return
	}

	utils.SuccessResponse(w, c, "OK")
}
