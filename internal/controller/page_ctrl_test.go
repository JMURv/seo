package ctrl

import (
	"context"
	"errors"
	"fmt"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/mocks"
	"github.com/JMURv/seo-svc/pkg/consts"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestController_ListPages(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	var expected []*model.Page

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().ListPages(gomock.Any()).Return(expected, nil).Times(1)

		res, err := ctrl.ListPages(ctx)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		newErr := errors.New("err")
		mockRepo.EXPECT().ListPages(gomock.Any()).Return(nil, newErr).Times(1)

		res, err := ctrl.ListPages(ctx)
		assert.IsType(t, newErr, err)
		assert.Nil(t, res)
	})
}

func TestController_GetPage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	slug := "slug"
	key := fmt.Sprintf(pageKey, slug)
	expected := &model.Page{}

	t.Run("Cache hit", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), key, gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, dest *model.Page) error {
				*dest = *expected
				return nil
			},
		).Times(1)

		res, err := ctrl.GetPage(ctx, slug)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Cache miss, repo success, cache set success", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetPage(gomock.Any(), slug).
			Return(expected, nil).
			Times(1)
		mockCache.EXPECT().
			Set(gomock.Any(), consts.DefaultCacheTime, key, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := ctrl.GetPage(ctx, slug)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Cache miss, repo returns ErrNotFound", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetPage(gomock.Any(), slug).
			Return(nil, repo.ErrNotFound).
			Times(1)

		res, err := ctrl.GetPage(ctx, slug)
		assert.Nil(t, res)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("Cache miss, repo error (other than ErrNotFound)", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetPage(gomock.Any(), slug).
			Return(nil, errors.New("some repo error")).
			Times(1)

		res, err := ctrl.GetPage(ctx, slug)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Cache miss, repo success, cache set failure", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetPage(gomock.Any(), slug).
			Return(expected, nil).
			Times(1)
		mockCache.EXPECT().
			Set(gomock.Any(), consts.DefaultCacheTime, key, gomock.Any()).
			Return(errors.New("cache set failure")).
			Times(1)

		res, err := ctrl.GetPage(ctx, slug)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestController_CreatePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	expected := "slug"
	req := &model.Page{
		Slug:  "slug",
		Title: "title",
		Href:  "href",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			CreatePage(gomock.Any(), req).
			Return(expected, nil).
			Times(1)

		res, err := ctrl.CreatePage(ctx, req)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		exp := ""
		mockRepo.EXPECT().
			CreatePage(gomock.Any(), req).
			Return(exp, repo.ErrAlreadyExists).
			Times(1)

		res, err := ctrl.CreatePage(ctx, req)
		assert.IsType(t, ErrAlreadyExists, err)
		assert.Equal(t, exp, res)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		exp := ""
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			CreatePage(gomock.Any(), req).
			Return(exp, newErr).
			Times(1)

		res, err := ctrl.CreatePage(ctx, req)
		assert.IsType(t, newErr, err)
		assert.Equal(t, exp, res)
	})
}

func TestController_UpdatePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	slug := "slug"
	req := &model.Page{
		Slug:  slug,
		Title: "title",
		Href:  "href",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			UpdatePage(gomock.Any(), slug, req).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(pageKey, slug)).
			Return(nil).
			Times(1)

		err := ctrl.UpdatePage(ctx, slug, req)
		assert.Nil(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			UpdatePage(gomock.Any(), slug, req).
			Return(repo.ErrNotFound).
			Times(1)

		err := ctrl.UpdatePage(ctx, slug, req)
		assert.IsType(t, ErrNotFound, err)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			UpdatePage(gomock.Any(), slug, req).
			Return(newErr).
			Times(1)

		err := ctrl.UpdatePage(ctx, slug, req)
		assert.IsType(t, newErr, err)
	})

	t.Run("ErrCache", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			UpdatePage(gomock.Any(), slug, req).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(pageKey, slug)).
			Return(newErr).
			Times(1)

		err := ctrl.UpdatePage(ctx, slug, req)
		assert.Nil(t, err)
	})
}

func TestController_DeletePage(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	slug := "slug"
	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			DeletePage(gomock.Any(), slug).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(pageKey, slug)).
			Return(nil).
			Times(1)

		err := ctrl.DeletePage(ctx, slug)
		assert.Nil(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			DeletePage(gomock.Any(), slug).
			Return(repo.ErrNotFound).
			Times(1)

		err := ctrl.DeletePage(ctx, slug)
		assert.IsType(t, ErrNotFound, err)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			DeletePage(gomock.Any(), slug).
			Return(newErr).
			Times(1)

		err := ctrl.DeletePage(ctx, slug)
		assert.IsType(t, newErr, err)
	})

	t.Run("ErrCache", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			DeletePage(gomock.Any(), slug).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(pageKey, slug)).
			Return(newErr).
			Times(1)

		err := ctrl.DeletePage(ctx, slug)
		assert.Nil(t, err)
	})
}
