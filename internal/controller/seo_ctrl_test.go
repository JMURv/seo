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

func TestController_GetSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	name, pk := "name", "pk"
	key := fmt.Sprintf(SEOKey, name, pk)
	expected := &model.SEO{}

	t.Run("Cache hit", func(t *testing.T) {
		mockCache.EXPECT().Get(gomock.Any(), key, gomock.Any()).DoAndReturn(
			func(_ context.Context, _ string, dest *model.SEO) error {
				*dest = *expected
				return nil
			},
		).Times(1)

		user, err := ctrl.GetSEO(ctx, name, pk)
		assert.Nil(t, err)
		assert.Equal(t, expected, user)
	})

	t.Run("Cache miss, repo success, cache set success", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(expected, nil).
			Times(1)
		mockCache.EXPECT().
			Set(gomock.Any(), consts.DefaultCacheTime, key, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := ctrl.GetSEO(ctx, name, pk)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("Cache miss, repo returns ErrNotFound", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(nil, repo.ErrNotFound).
			Times(1)

		res, err := ctrl.GetSEO(ctx, name, pk)
		assert.Nil(t, res)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("Cache miss, repo error (other than ErrNotFound)", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(nil, errors.New("some repo error")).
			Times(1)

		res, err := ctrl.GetSEO(ctx, name, pk)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Cache miss, repo success, cache set failure", func(t *testing.T) {
		mockCache.EXPECT().
			Get(gomock.Any(), key, gomock.Any()).
			Return(errors.New("cache miss")).
			Times(1)
		mockRepo.EXPECT().
			GetSEO(gomock.Any(), name, pk).
			Return(expected, nil).
			Times(1)
		mockCache.EXPECT().
			Set(gomock.Any(), consts.DefaultCacheTime, key, gomock.Any()).
			Return(errors.New("cache set failure")).
			Times(1)

		res, err := ctrl.GetSEO(ctx, name, pk)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})
}

func TestController_CreateSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	name, pk := "name", "pk"
	expected := uint64(1)

	req := &model.SEO{
		Title:         "title",
		Description:   "description",
		Keywords:      "keyword1, keyword2",
		OGTitle:       "OG title",
		OGDescription: "OG description",
		OGImage:       "OG image",
		OBJName:       name,
		OBJPK:         pk,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			CreateSEO(gomock.Any(), req).
			Return(uint64(1), nil).
			Times(1)

		res, err := ctrl.CreateSEO(ctx, req)
		assert.Nil(t, err)
		assert.Equal(t, expected, res)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		exp := uint64(0)
		mockRepo.EXPECT().
			CreateSEO(gomock.Any(), req).
			Return(exp, repo.ErrAlreadyExists).
			Times(1)

		res, err := ctrl.CreateSEO(ctx, req)
		assert.IsType(t, ErrAlreadyExists, err)
		assert.Equal(t, exp, res)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		exp := uint64(0)
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			CreateSEO(gomock.Any(), req).
			Return(exp, newErr).
			Times(1)

		res, err := ctrl.CreateSEO(ctx, req)
		assert.IsType(t, newErr, err)
		assert.Equal(t, exp, res)
	})
}

func TestController_UpdateSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	name, pk := "name", "pk"
	req := &model.SEO{
		Title:         "title",
		Description:   "description",
		Keywords:      "keyword1, keyword2",
		OGTitle:       "OG title",
		OGDescription: "OG description",
		OGImage:       "OG image",
		OBJName:       name,
		OBJPK:         pk,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			UpdateSEO(gomock.Any(), req).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(SEOKey, name, pk)).
			Return(nil).
			Times(1)

		err := ctrl.UpdateSEO(ctx, req)
		assert.Nil(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			UpdateSEO(gomock.Any(), req).
			Return(repo.ErrNotFound).
			Times(1)

		err := ctrl.UpdateSEO(ctx, req)
		assert.IsType(t, ErrNotFound, err)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			UpdateSEO(gomock.Any(), req).
			Return(newErr).
			Times(1)

		err := ctrl.UpdateSEO(ctx, req)
		assert.IsType(t, newErr, err)
	})

	t.Run("ErrCache", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			UpdateSEO(gomock.Any(), req).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(SEOKey, name, pk)).
			Return(newErr).
			Times(1)

		err := ctrl.UpdateSEO(ctx, req)
		assert.Nil(t, err)
	})
}

func TestController_DeleteSEO(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockappRepo(ctrlMock)
	mockCache := mocks.NewMockCacheRepo(ctrlMock)

	ctx := context.Background()
	ctrl := New(mockRepo, mockCache)

	name, pk := "name", "pk"
	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(SEOKey, name, pk)).
			Return(nil).
			Times(1)

		err := ctrl.DeleteSEO(ctx, name, pk)
		assert.Nil(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockRepo.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(repo.ErrNotFound).
			Times(1)

		err := ctrl.DeleteSEO(ctx, name, pk)
		assert.IsType(t, ErrNotFound, err)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(newErr).
			Times(1)

		err := ctrl.DeleteSEO(ctx, name, pk)
		assert.IsType(t, newErr, err)
	})

	t.Run("ErrCache", func(t *testing.T) {
		newErr := errors.New("some error")
		mockRepo.EXPECT().
			DeleteSEO(gomock.Any(), name, pk).
			Return(nil).
			Times(1)
		mockCache.EXPECT().
			Delete(gomock.Any(), fmt.Sprintf(SEOKey, name, pk)).
			Return(newErr).
			Times(1)

		err := ctrl.DeleteSEO(ctx, name, pk)
		assert.Nil(t, err)
	})
}
