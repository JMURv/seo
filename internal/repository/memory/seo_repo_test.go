package memory

import (
	"context"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRepository_GetSEO(t *testing.T) {
	r := New()
	ctx := context.Background()

	seoEntry := &model.SEO{
		ID:          1,
		OBJName:     "test-name",
		OBJPK:       "test-pk",
		Title:       "Test Title",
		Description: "Test Description",
		Keywords:    "Test Keywords",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.SEOData[seoEntry.ID] = seoEntry

	t.Run("Success", func(t *testing.T) {
		res, err := r.GetSEO(ctx, "test-name", "test-pk")
		require.NoError(t, err)
		assert.Equal(t, seoEntry, res)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		res, err := r.GetSEO(ctx, "invalid-name", "invalid-pk")
		assert.ErrorIs(t, err, repo.ErrNotFound)
		assert.Nil(t, res)
	})
}

func TestRepository_CreateSEO(t *testing.T) {
	r := New()
	ctx := context.Background()

	existingSEO := &model.SEO{
		OBJName:     "duplicate-name",
		OBJPK:       "duplicate-pk",
		Title:       "Duplicate Title",
		Description: "Duplicate Description",
		Keywords:    "Keyword1, Keyword2",
	}
	_, err := r.CreateSEO(ctx, existingSEO)
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		newSEO := &model.SEO{
			OBJName:     "new-name",
			OBJPK:       "new-pk",
			Title:       "New Title",
			Description: "New Description",
			Keywords:    "Keyword1, Keyword2",
		}
		_, err := r.CreateSEO(ctx, newSEO)
		require.NoError(t, err)

		createdSEO, err := r.GetSEO(ctx, newSEO.OBJName, newSEO.OBJPK)
		require.NoError(t, err)
		assert.Equal(t, newSEO, createdSEO)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		_, err = r.CreateSEO(ctx, existingSEO)
		assert.ErrorIs(t, err, repo.ErrAlreadyExists)
	})
}

func TestRepository_UpdateSEO(t *testing.T) {
	r := New()
	ctx := context.Background()
	existingSEO := &model.SEO{
		ID:          1,
		OBJName:     "test-name",
		OBJPK:       "test-pk",
		Title:       "Original Title",
		Description: "Original Description",
		Keywords:    "Keyword1",
	}
	r.SEOData[existingSEO.ID] = existingSEO

	t.Run("Success", func(t *testing.T) {
		updatedSEO := &model.SEO{
			ID:          1,
			OBJName:     "test-name",
			OBJPK:       "test-pk",
			Title:       "Updated Title",
			Description: "Updated Description",
			Keywords:    "Keyword1, Keyword2",
		}

		err := r.UpdateSEO(ctx, updatedSEO)
		require.NoError(t, err)

		updatedEntry, err := r.GetSEO(ctx, updatedSEO.OBJName, updatedSEO.OBJPK)
		require.NoError(t, err)
		assert.Equal(t, updatedSEO, updatedEntry)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		nonExistingSEO := &model.SEO{
			ID:      999,
			OBJName: "nonexistent-name",
			OBJPK:   "nonexistent-pk",
			Title:   "Nonexistent Title",
		}

		err := r.UpdateSEO(context.Background(), nonExistingSEO)
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}

func TestRepository_DeleteSEO(t *testing.T) {
	r := New()
	ctx := context.Background()
	existingSEO := &model.SEO{
		ID:          1,
		OBJName:     "test-name",
		OBJPK:       "test-pk",
		Title:       "To Be Deleted",
		Description: "To Be Deleted",
	}
	r.SEOData[existingSEO.ID] = existingSEO

	t.Run("Success", func(t *testing.T) {
		err := r.DeleteSEO(ctx, "test-name", "test-pk")
		require.NoError(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		err := r.DeleteSEO(context.Background(), "nonexistent-name", "nonexistent-pk")
		assert.ErrorIs(t, err, repo.ErrNotFound)
	})
}
