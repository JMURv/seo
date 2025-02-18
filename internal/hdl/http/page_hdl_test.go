package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	ctrl "github.com/JMURv/seo/internal/ctrl"
	model "github.com/JMURv/seo/internal/models"
	"github.com/JMURv/seo/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ListPages(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/page"
	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	ctx := context.Background()

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				ListPages(gomock.Any()).
				Return([]*model.Page{}, nil).
				Times(1)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.ListPages(w, req)
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrInternalError", func(t *testing.T) {
			var ErrOther = errors.New("other error")
			mockCtrl.EXPECT().
				ListPages(gomock.Any()).
				Return(nil, ErrOther).
				Times(1)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.ListPages(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		},
	)
}

func TestHandler_GetPage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/page/slug"
	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	ctx := context.Background()
	slug := "slug"

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				GetPage(gomock.Any(), slug).
				Return(&model.Page{}, nil).
				Times(1)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.GetPage(w, req)
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		},
	)

	t.Run(
		"Missing slug", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/page/", nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.GetPage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				GetPage(gomock.Any(), slug).
				Return(nil, ctrl.ErrNotFound).
				Times(1)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.GetPage(w, req)
			assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrInternalError", func(t *testing.T) {
			var ErrOther = errors.New("other error")
			mockCtrl.EXPECT().
				GetPage(gomock.Any(), slug).
				Return(nil, ErrOther).
				Times(1)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.GetPage(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		},
	)
}

func TestHandler_CreatePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	slug := "slug"
	ctx := context.Background()

	reqData := &model.Page{
		Slug:  "slug",
		Title: "name",
		Href:  "href",
	}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				CreatePage(gomock.Any(), reqData).
				Return(slug, nil).
				Times(1)

			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPost, "/api/page", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreatePage(w, req)
			assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrAlreadyExists", func(t *testing.T) {
			mockCtrl.EXPECT().
				CreatePage(gomock.Any(), reqData).
				Return("", ctrl.ErrAlreadyExists).
				Times(1)

			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPost, "/api/page", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreatePage(w, req)
			assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrInternalError", func(t *testing.T) {
			var ErrOther = errors.New("other error")
			mockCtrl.EXPECT().
				CreatePage(gomock.Any(), reqData).
				Return("", ErrOther).
				Times(1)

			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPost, "/api/page", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreatePage(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrDecodeRequest - Missing Title", func(t *testing.T) {
			reqData.Title = ""
			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPost, "/api/page", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreatePage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)

	t.Run(
		"Invalid JSON", func(t *testing.T) {
			payload, _ := json.Marshal(map[string]any{"title": 123})
			req := httptest.NewRequest(http.MethodPost, "/api/page", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.CreatePage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)
}

func TestHandler_UpdatePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/page/slug"
	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	slug := "slug"
	ctx := context.Background()

	reqData := &model.Page{
		Slug:  "slug",
		Title: "name",
		Href:  "href",
	}

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				UpdatePage(gomock.Any(), slug, reqData).
				Return(nil).
				Times(1)

			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdatePage(w, req)
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		},
	)

	t.Run(
		"Missing name or pk", func(t *testing.T) {
			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPut, "/api/page/", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdatePage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				UpdatePage(gomock.Any(), slug, reqData).
				Return(ctrl.ErrNotFound).
				Times(1)

			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdatePage(w, req)
			assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrInternalError", func(t *testing.T) {
			var ErrOther = errors.New("other error")
			mockCtrl.EXPECT().
				UpdatePage(gomock.Any(), slug, reqData).
				Return(ErrOther).
				Times(1)

			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdatePage(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrDecodeRequest - Missing Title", func(t *testing.T) {
			reqData.Title = ""
			payload, _ := json.Marshal(reqData)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdatePage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)

	t.Run(
		"Invalid JSON", func(t *testing.T) {
			payload, _ := json.Marshal(map[string]any{"title": 123})
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.UpdatePage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)
}

func TestHandler_DeletePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/page/slug"
	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	slug := "slug"
	ctx := context.Background()

	t.Run(
		"Success", func(t *testing.T) {
			mockCtrl.EXPECT().
				DeletePage(gomock.Any(), slug).
				Return(nil).
				Times(1)

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeletePage(w, req)
			assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
		},
	)

	t.Run(
		"Missing name or pk", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/seo/test-name/", nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeletePage(w, req)
			assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mockCtrl.EXPECT().
				DeletePage(gomock.Any(), slug).
				Return(ctrl.ErrNotFound).
				Times(1)

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeletePage(w, req)
			assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		},
	)

	t.Run(
		"ErrInternalError", func(t *testing.T) {
			var ErrOther = errors.New("other error")
			mockCtrl.EXPECT().
				DeletePage(gomock.Any(), slug).
				Return(ErrOther).
				Times(1)

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			h.DeletePage(w, req)
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		},
	)
}
