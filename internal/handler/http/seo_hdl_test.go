package http

import (
	"bytes"
	"context"
	"errors"
	ctrl "github.com/JMURv/seo-svc/internal/controller"
	"github.com/JMURv/seo-svc/mocks"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/seo/name/pk"
	mockCtrl := mocks.NewMockSEOCtrl(ctrlMock)
	h := New(mockCtrl)

	ctx := context.Background()
	name, pk := "name", "pk"

	t.Run("Success", func(t *testing.T) {
		mockCtrl.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(&model.SEO{}, nil).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetSEO(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Missing name or pk", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/seo/test-name/", nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockCtrl.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(nil, ctrl.ErrNotFound).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetSEO(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		mockCtrl.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(nil, ErrOther).
			Times(1)

		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.GetSEO(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestHandler_CreateSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockSEOCtrl(ctrlMock)
	h := New(mockCtrl)

	ctx := context.Background()

	reqData := &model.SEO{
		Title:         "name",
		Description:   "description",
		Keywords:      "keywords",
		OGTitle:       "ogtitle",
		OGDescription: "ogdescription",
		OGImage:       "ogimage",
		OBJName:       "objname",
		OBJPK:         "objpk",
	}

	t.Run("Success", func(t *testing.T) {
		mockCtrl.EXPECT().
			CreateSEO(gomock.Any(), reqData).
			Return(uint64(1), nil).
			Times(1)

		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPost, "/api/seo", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.CreateSEO(w, req)
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		mockCtrl.EXPECT().
			CreateSEO(gomock.Any(), reqData).
			Return(uint64(0), ctrl.ErrAlreadyExists).
			Times(1)

		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPost, "/api/seo", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.CreateSEO(w, req)
		assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		mockCtrl.EXPECT().
			CreateSEO(gomock.Any(), reqData).
			Return(uint64(0), ErrOther).
			Times(1)

		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPost, "/api/seo", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.CreateSEO(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("ErrDecodeRequest - Missing Title", func(t *testing.T) {
		reqData.Title = ""
		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPost, "/api/seo", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.CreateSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]any{"title": 123})
		req := httptest.NewRequest(http.MethodPost, "/api/seo", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.CreateSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_UpdateSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/seo/test-name/1"
	mockCtrl := mocks.NewMockSEOCtrl(ctrlMock)
	h := New(mockCtrl)

	ctx := context.Background()

	reqData := &model.SEO{
		Title:         "name",
		Description:   "description",
		Keywords:      "keywords",
		OGTitle:       "ogtitle",
		OGDescription: "ogdescription",
		OGImage:       "ogimage",
		OBJName:       "test-name",
		OBJPK:         "1",
	}

	t.Run("Success", func(t *testing.T) {
		mockCtrl.EXPECT().
			UpdateSEO(gomock.Any(), reqData).
			Return(nil).
			Times(1)

		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.UpdateSEO(w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Missing name or pk", func(t *testing.T) {
		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPut, "/api/seo/test-name/", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.UpdateSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockCtrl.EXPECT().
			UpdateSEO(gomock.Any(), reqData).
			Return(ctrl.ErrNotFound).
			Times(1)

		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.UpdateSEO(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		mockCtrl.EXPECT().
			UpdateSEO(gomock.Any(), reqData).
			Return(ErrOther).
			Times(1)

		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.UpdateSEO(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("ErrDecodeRequest - Missing Title", func(t *testing.T) {
		reqData.Title = ""
		payload, _ := json.Marshal(reqData)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.UpdateSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		payload, _ := json.Marshal(map[string]any{"title": 123})
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.UpdateSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestHandler_DeleteSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/seo/name/pk"
	mockCtrl := mocks.NewMockSEOCtrl(ctrlMock)
	h := New(mockCtrl)

	ctx := context.Background()
	name, pk := "name", "pk"

	t.Run("Success", func(t *testing.T) {
		mockCtrl.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(nil).
			Times(1)

		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.DeleteSEO(w, req)
		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("Missing name or pk", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/seo/test-name/", nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.DeleteSEO(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockCtrl.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(ctrl.ErrNotFound).
			Times(1)

		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.DeleteSEO(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("ErrInternalError", func(t *testing.T) {
		var ErrOther = errors.New("other error")
		mockCtrl.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(ErrOther).
			Times(1)

		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		h.DeleteSEO(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}
