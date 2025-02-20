package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/JMURv/seo/internal/dto"
	model "github.com/JMURv/seo/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"testing"
)

func TestPages(t *testing.T) {
	server, authReq, cleanup := setupTestServer()
	t.Cleanup(cleanup)
	t.Cleanup(server.Close)

	tokenRes := authReq(context.Background(), "admin@example.com", "superstrongpassword")
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", tokenRes),
	}

	listPages := func() []model.Page {
		resp, err := http.Get(server.URL + "/api/page")
		require.Nil(t, err)

		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				zap.L().Debug("failed to close response body", zap.Error(err))
			}
		}(resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var pages []model.Page
		require.Nil(t, json.NewDecoder(resp.Body).Decode(&pages))
		return pages
	}

	getPage := func(slug string) *model.Page {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/api/page/"+slug, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
				zap.L().Debug("failed to close response body", zap.Error(err))
			}
		}(resp.Body)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var r model.Page
		require.Nil(t, json.NewDecoder(resp.Body).Decode(&r))
		return &r
	}

	createPage := func(page *model.Page, headers map[string]string) *dto.CreatePageResponse {
		payload, err := json.Marshal(page)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/api/page", bytes.NewBuffer(payload))
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

		var r dto.CreatePageResponse
		require.Nil(t, json.NewDecoder(resp.Body).Decode(&r))
		return &r
	}

	updatePage := func(slug string, page *model.Page, headers map[string]string) {
		payload, err := json.Marshal(page)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPut, server.URL+"/api/page/"+slug, bytes.NewBuffer(payload))
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

	deletePage := func(slug string, headers map[string]string) {
		req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/page/"+slug, nil)
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

	assert.Equal(t, 0, len(listPages()))

	pageReq := &model.Page{
		Slug:  "slug",
		Title: "title",
		Href:  "href",
	}
	newlyPage := createPage(pageReq, headers)
	assert.Equal(t, 1, len(listPages()))

	page := getPage(newlyPage.Slug)
	assert.Equal(t, pageReq.Slug, page.Slug)
	assert.Equal(t, pageReq.Title, page.Title)
	assert.Equal(t, pageReq.Href, page.Href)

	pageReq.Title = "new title"
	updatePage(newlyPage.Slug, pageReq, headers)

	page = getPage(newlyPage.Slug)
	assert.Equal(t, pageReq.Slug, page.Slug)
	assert.Equal(t, pageReq.Title, page.Title)
	assert.Equal(t, pageReq.Href, page.Href)

	deletePage(newlyPage.Slug, headers)
	assert.Equal(t, 0, len(listPages()))
}
