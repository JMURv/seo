package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/JMURv/seo/internal/ctrl"
	"github.com/JMURv/seo/internal/dto"
	md "github.com/JMURv/seo/internal/models"
	"github.com/JMURv/seo/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetSEO(t *testing.T) {
	cmock := gomock.NewController(t)
	defer cmock.Finish()

	const url = "/api/seo/name/pk"
	mctrl := mocks.NewMockAppCtrl(cmock)
	sso := mocks.NewMockSSOSvc(cmock)
	h := New(mctrl, sso)

	ctx := context.Background()
	name, pk := "name", "pk"
	testErr := errors.New("test error")

	tests := []struct {
		name   string
		url    string
		method string
		status int
		expect func()
	}{
		{
			name:   "Success",
			url:    url,
			method: http.MethodGet,
			status: http.StatusOK,
			expect: func() {
				mctrl.EXPECT().
					GetSEO(gomock.Any(), name, pk).
					Return(&md.SEO{}, nil).
					Times(1)
			},
		},
		{
			name:   "Missing name or pk",
			url:    "/api/seo/test-name/",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			expect: func() {},
		},
		{
			name:   "ErrNotFound",
			url:    url,
			method: http.MethodGet,
			status: http.StatusNotFound,
			expect: func() {
				mctrl.EXPECT().
					GetSEO(gomock.Any(), name, pk).
					Return(nil, ctrl.ErrNotFound).
					Times(1)
			},
		},
		{
			name:   "ErrInternal",
			url:    url,
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					GetSEO(gomock.Any(), name, pk).
					Return(nil, testErr).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.expect()
				req := httptest.NewRequestWithContext(ctx, tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				h.GetSEO(w, req)
				assert.Equal(t, tt.status, w.Result().StatusCode)
			},
		)
	}
}

func TestHandler_CreateSEO(t *testing.T) {
	const url = "/api/seo"
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	ctx := context.Background()
	var ErrOther = errors.New("other error")

	reqData := &md.SEO{
		Title:         "name",
		Description:   "description",
		Keywords:      "keywords",
		OGTitle:       "ogtitle",
		OGDescription: "ogdescription",
		OGImage:       "ogimage",
		OBJName:       "objname",
		OBJPK:         "objpk",
	}

	tests := []struct {
		ctx     context.Context
		name    string
		url     string
		method  string
		status  int
		payload map[string]any
		expect  func()
	}{
		{
			name:   "ErrInternal",
			url:    url,
			method: http.MethodPost,
			status: http.StatusInternalServerError,
			payload: map[string]any{
				"title":         reqData.Title,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {
				mockCtrl.EXPECT().
					CreateSEO(gomock.Any(), reqData).
					Return(nil, ErrOther).
					Times(1)
			},
		},
		{
			name:   "MissingTitle",
			url:    url,
			method: http.MethodPost,
			status: http.StatusBadRequest,
			payload: map[string]any{
				"title":         "",
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {},
		},
		{
			name:   "ErrDecodeRequest",
			url:    url,
			method: http.MethodPost,
			status: http.StatusBadRequest,
			payload: map[string]any{
				"title":         0,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {},
		},
		{
			name:   "ErrAlreadyExists",
			url:    url,
			method: http.MethodPost,
			status: http.StatusConflict,
			payload: map[string]any{
				"title":         reqData.Title,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {
				mockCtrl.EXPECT().
					CreateSEO(gomock.Any(), reqData).
					Return(nil, ctrl.ErrAlreadyExists).
					Times(1)
			},
		},
		{
			name:   "Success",
			url:    url,
			method: http.MethodPost,
			status: http.StatusCreated,
			payload: map[string]any{
				"title":         reqData.Title,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {
				mockCtrl.EXPECT().
					CreateSEO(gomock.Any(), reqData).
					Return(
						&dto.CreateSEOResponse{
							Name: "name",
							PK:   "pk",
						}, nil,
					).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.expect()

				payload, err := json.Marshal(tt.payload)
				assert.Nil(t, err)

				req := httptest.NewRequestWithContext(ctx, tt.method, tt.url, bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				h.CreateSEO(w, req)
				assert.Equal(t, tt.status, w.Result().StatusCode)
			},
		)
	}
}

func TestHandler_UpdateSEO(t *testing.T) {
	const url = "/api/seo/test-name/1"
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	ctx := context.Background()
	var ErrOther = errors.New("other error")

	reqData := &md.SEO{
		Title:         "name",
		Description:   "description",
		Keywords:      "keywords",
		OGTitle:       "ogtitle",
		OGDescription: "OGDescription",
		OGImage:       "ogimage",
		OBJName:       "test-name",
		OBJPK:         "1",
	}

	tests := []struct {
		ctx     context.Context
		name    string
		url     string
		method  string
		status  int
		payload map[string]any
		expect  func()
	}{
		{
			name:   "Success",
			url:    url,
			method: http.MethodPut,
			status: http.StatusOK,
			payload: map[string]any{
				"title":         reqData.Title,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {
				mockCtrl.EXPECT().
					UpdateSEO(gomock.Any(), reqData).
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "Missing name or pk",
			url:    "/api/seo/test-name/",
			method: http.MethodPut,
			status: http.StatusBadRequest,
			payload: map[string]any{
				"title": reqData.Title,
			},
			expect: func() {},
		},
		{
			name:   "ErrNotFound",
			url:    url,
			method: http.MethodPut,
			status: http.StatusNotFound,
			payload: map[string]any{
				"title":         reqData.Title,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {
				mockCtrl.EXPECT().
					UpdateSEO(gomock.Any(), reqData).
					Return(ctrl.ErrNotFound).
					Times(1)
			},
		},
		{
			name:   "ErrInternal",
			url:    url,
			method: http.MethodPut,
			status: http.StatusInternalServerError,
			payload: map[string]any{
				"title":         reqData.Title,
				"description":   reqData.Description,
				"keywords":      reqData.Keywords,
				"OGTitle":       reqData.OGTitle,
				"OGDescription": reqData.OGDescription,
				"OGImage":       reqData.OGImage,
				"obj_name":      reqData.OBJName,
				"obj_pk":        reqData.OBJPK,
			},
			expect: func() {
				mockCtrl.EXPECT().
					UpdateSEO(gomock.Any(), reqData).
					Return(ErrOther).
					Times(1)
			},
		},
		{
			name:   "Validation Error",
			url:    url,
			method: http.MethodPut,
			status: http.StatusBadRequest,
			payload: map[string]any{
				"title": "",
			},
			expect: func() {},
		},
		{
			name:   "ErrDecodeRequest",
			url:    url,
			method: http.MethodPut,
			status: http.StatusBadRequest,
			payload: map[string]any{
				"title": 123,
			},
			expect: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.expect()

				payload, err := json.Marshal(tt.payload)
				assert.Nil(t, err)

				req := httptest.NewRequestWithContext(ctx, tt.method, tt.url, bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				h.UpdateSEO(w, req)
				assert.Equal(t, tt.status, w.Result().StatusCode)
			},
		)
	}
}

func TestHandler_DeleteSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	const url = "/api/seo/name/pk"
	mockCtrl := mocks.NewMockAppCtrl(ctrlMock)
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)
	h := New(mockCtrl, ssoCtrl)

	ctx := context.Background()
	name, pk := "name", "pk"
	ErrOther := errors.New("other error")

	tests := []struct {
		ctx    context.Context
		name   string
		url    string
		method string
		status int
		expect func()
	}{
		{
			name:   "Success",
			url:    url,
			method: http.MethodDelete,
			status: http.StatusNoContent,
			expect: func() {
				mockCtrl.EXPECT().
					DeleteSEO(gomock.Any(), name, pk).
					Return(nil).
					Times(1)
			},
		},
		{
			name:   "Missing name or pk",
			url:    "/api/seo/test-name/",
			method: http.MethodDelete,
			status: http.StatusBadRequest,
			expect: func() {},
		},
		{
			name:   "ErrNotFound",
			url:    url,
			method: http.MethodDelete,
			status: http.StatusNotFound,
			expect: func() {
				mockCtrl.EXPECT().
					DeleteSEO(gomock.Any(), name, pk).
					Return(ctrl.ErrNotFound).
					Times(1)
			},
		},
		{
			name:   "ErrInternal",
			url:    url,
			method: http.MethodDelete,
			status: http.StatusInternalServerError,
			expect: func() {
				mockCtrl.EXPECT().
					DeleteSEO(gomock.Any(), name, pk).
					Return(ErrOther).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				tt.expect()

				req := httptest.NewRequestWithContext(ctx, tt.method, tt.url, nil)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				h.DeleteSEO(w, req)
				assert.Equal(t, tt.status, w.Result().StatusCode)
			},
		)
	}
}
