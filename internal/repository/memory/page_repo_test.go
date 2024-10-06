package memory

import (
	"context"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRepository_GetPage(t *testing.T) {
	r := New()
	ctx := context.Background()

	pageEntry := &model.Page{
		Slug:  "test-slug",
		Title: "Test Title",
		Href:  "/test-slug",
	}
	r.PageData[pageEntry.Slug] = pageEntry

	t.Run("Success", func(t *testing.T) {
		res, err := r.GetPage(ctx, "test-slug")
		require.NoError(t, err)
		assert.Equal(t, pageEntry, res)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		res, err := r.GetPage(ctx, "nonexistent-slug")
		assert.ErrorIs(t, err, repo.ErrNotFound)
		assert.Nil(t, res)
	})
}

func TestRepository_CreatePage(t *testing.T) {
	r := New()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		newPage := &model.Page{
			Slug:  "new-slug",
			Title: "New Title",
			Href:  "/new-slug",
		}

		slug, err := r.CreatePage(ctx, newPage)
		require.NoError(t, err)
		assert.Equal(t, "new-slug", slug)

		createdPage, err := r.GetPage(ctx, newPage.Slug)
		require.NoError(t, err)
		assert.Equal(t, newPage, createdPage)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		existingPage := &model.Page{
			Slug:  "duplicate-slug",
			Title: "Duplicate Title",
			Href:  "/duplicate-slug",
		}
		_, err := r.CreatePage(ctx, existingPage)
		require.NoError(t, err)

		_, err = r.CreatePage(ctx, existingPage)
		assert.ErrorIs(t, err, repo.ErrAlreadyExists)
	})
}

func TestRepository_UpdatePage(t *testing.T) {
	r := New()
	ctx := context.Background()

	existingPage := &model.Page{
		Slug:  "test-slug",
		Title: "Original Title",
		Href:  "/original-slug",
	}
	r.PageData[existingPage.Slug] = existingPage

	t.Run("Success", func(t *testing.T) {
		updatedPage := &model.Page{
			Slug:  "test-slug",
			Title: "Updated Title",
			Href:  "/updated-slug",
		}

		err := r.UpdatePage(ctx, "test-slug", updatedPage)
		require.NoError(t, err)

		updatedEntry, err := r.GetPage(ctx, updatedPage.Slug)
		require.NoError(t, err)
		assert.Equal(t, updatedPage, updatedEntry)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		nonExistingPage := &model.Page{
			Slug:  "nonexistent-slug",
			Title: "Nonexistent Title",
		}

		err := r.UpdatePage(ctx, "nonexistent-slug", nonExistingPage)
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}

func TestRepository_DeletePage(t *testing.T) {
	r := New()
	ctx := context.Background()

	existingPage := &model.Page{
		Slug:  "test-slug",
		Title: "To Be Deleted",
		Href:  "/to-be-deleted",
	}
	r.PageData[existingPage.Slug] = existingPage

	t.Run("Success", func(t *testing.T) {
		err := r.DeletePage(ctx, "test-slug")
		require.NoError(t, err)

		_, err = r.GetPage(ctx, "test-slug")
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		err := r.DeletePage(ctx, "nonexistent-slug")
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}
