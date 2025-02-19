package http

import (
	"bytes"
	"encoding/json"
	model "github.com/JMURv/seo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"
)

func TestSEO(t *testing.T) {
	server, sso, cleanup := setupTestServer(t)
	t.Cleanup(cleanup)
	t.Cleanup(server.Close)

	name, pk := "page", "slug"
	headers := map[string]string{
		"Authorization": "Bearer valid-token",
	}

	getSEO := func(objName, objPK string) *model.SEO {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/seo/"+objName+"/"+objPK, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				zap.L().Debug("failed to close response body", zap.Error(err))
			}
		}(resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var r model.SEO
		require.Nil(t, json.NewDecoder(resp.Body).Decode(&r))
		return &r
	}

	createSEO := func(seo *model.SEO, headers map[string]string) {
		sso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return("token", nil)

		payload, err := json.Marshal(seo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/seo", bytes.NewBuffer(payload))
		require.NoError(t, err)

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				zap.L().Debug("failed to close response body", zap.Error(err))
			}
		}(resp.Body)

		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	updateSEO := func(objName, objPK string, seo *model.SEO, headers map[string]string) {
		sso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return("token", nil)

		payload, err := json.Marshal(seo)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/seo/"+objName+"/"+objPK, bytes.NewBuffer(payload))
		require.NoError(t, err)

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				zap.L().Debug("failed to close response body", zap.Error(err))
			}
		}(resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)
	}

	deleteSEO := func(objName, objPK string, headers map[string]string) {
		sso.EXPECT().ParseClaims(gomock.Any(), gomock.Any()).Return("token", nil)

		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/seo/"+objName+"/"+objPK, nil)
		require.NoError(t, err)

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				zap.L().Debug("failed to close response body", zap.Error(err))
			}
		}(resp.Body)

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	}

	seoReq := &model.SEO{
		Title:         "title",
		Description:   "desc",
		Keywords:      "kw, kw1",
		OGTitle:       "title",
		OGDescription: "desc",
		OGImage:       "img/path",
		OBJName:       name,
		OBJPK:         pk,
	}
	createSEO(seoReq, headers)

	seo := getSEO(seoReq.OBJName, seoReq.OBJPK)
	assert.Equal(t, seoReq.Title, seo.Title)

	seoReq.Title = "new title"
	updateSEO(seoReq.OBJName, seoReq.OBJPK, seoReq, headers)

	seo = getSEO(seoReq.OBJName, seoReq.OBJPK)
	assert.Equal(t, seoReq.Title, seo.Title)

	deleteSEO(seoReq.OBJName, seoReq.OBJPK, headers)
}
