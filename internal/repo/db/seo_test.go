package db

import (
	"context"
	"database/sql"
	"errors"
	model "github.com/JMURv/seo/internal/models"
	rrepo "github.com/JMURv/seo/internal/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
)

func TestRepository_GetSEO(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	ctx := context.Background()

	name, pk := "name", "pk"
	testOBJ := model.SEO{
		Title:         "title",
		Description:   "description",
		Keywords:      "keywords1, keywords2",
		OGTitle:       "OGTitle",
		OGDescription: "OGDescription",
		OGImage:       "OGImage",
		OBJName:       name,
		OBJPK:         pk,
	}

	t.Run(
		"Success case", func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(getSEO)).
				WithArgs(name, pk).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{
							"title",
							"description",
							"keywords",
							"og_title",
							"og_description",
							"og_image",
							"obj_name",
							"obj_pk",
							"created_at",
							"updated_at",
						},
					).AddRow(
						testOBJ.Title,
						testOBJ.Description,
						testOBJ.Keywords,
						testOBJ.OGTitle,
						testOBJ.OGDescription,
						testOBJ.OGImage,
						testOBJ.OBJName,
						testOBJ.OBJPK,
						testOBJ.CreatedAt,
						testOBJ.UpdatedAt,
					),
				)

			result, err := repo.GetSEO(ctx, name, pk)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, testOBJ.ID, result.ID)
			assert.Equal(t, testOBJ.Title, result.Title)
			assert.Equal(t, testOBJ.Description, result.Description)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(getSEO)).
				WithArgs(name, pk).
				WillReturnError(sql.ErrNoRows)

			result, err := repo.GetSEO(ctx, name, pk)
			assert.Nil(t, result)
			assert.Equal(t, rrepo.ErrNotFound, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"Unexpected error case", func(t *testing.T) {
			notExpectedError := errors.New("not expected error")

			mock.ExpectQuery(regexp.QuoteMeta(getSEO)).
				WithArgs(name, pk).
				WillReturnError(notExpectedError)

			result, err := repo.GetSEO(ctx, name, pk)
			assert.Nil(t, result)
			assert.Equal(t, notExpectedError, err)
			assert.NotEqual(t, rrepo.ErrNotFound, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)
}

func TestRepository_CreateSEO(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}

	ctx := context.Background()
	name, pk := "name", "pk"
	testErr := errors.New("test error")
	testOBJ := &model.SEO{
		ID:            1,
		Title:         "title",
		Description:   "description",
		Keywords:      "keywords1, keywords2",
		OGTitle:       "OGTitle",
		OGDescription: "OGDescription",
		OGImage:       "OGImage",
		OBJName:       name,
		OBJPK:         pk,
	}

	t.Run(
		"Success case", func(t *testing.T) {
			mock.ExpectQuery(
				regexp.QuoteMeta(createSEO),
			).WillReturnRows(
				sqlmock.NewRows([]string{"obj_name", "obj_pk"}).
					AddRow(testOBJ.OBJName, testOBJ.OBJPK),
			)

			_, _, err := repo.CreateSEO(ctx, testOBJ)
			assert.NoError(t, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrInternal", func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(createSEO)).
				WillReturnError(testErr)

			_, _, err := repo.CreateSEO(ctx, testOBJ)
			assert.Equal(t, testErr, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)
}

func TestRepository_UpdateSEO(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	id, name, pk := 1, "name", "pk"
	testOBJ := &model.SEO{
		ID:            uint64(id),
		Title:         "title",
		Description:   "description",
		Keywords:      "keywords1, keywords2",
		OGTitle:       "OGTitle",
		OGDescription: "OGDescription",
		OGImage:       "OGImage",
		OBJName:       name,
		OBJPK:         pk,
	}

	mock.ExpectQuery(regexp.QuoteMeta(createSEO)).
		WillReturnRows(sqlmock.NewRows([]string{"obj_name", "obj_pk"}).AddRow(testOBJ.OBJName, testOBJ.OBJPK))

	_, _, err = repo.CreateSEO(context.Background(), testOBJ)
	require.NoError(t, err)

	t.Run(
		"Success", func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(updateSEO)).
				WithArgs(
					testOBJ.Title,
					testOBJ.Description,
					testOBJ.Keywords,
					testOBJ.OGTitle,
					testOBJ.OGDescription,
					testOBJ.OGImage,
					testOBJ.OBJName,
					testOBJ.OBJPK,
					testOBJ.OBJName,
					testOBJ.OBJPK,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := repo.UpdateSEO(context.Background(), testOBJ)
			assert.NoError(t, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(updateSEO)).
				WithArgs(
					testOBJ.Title,
					testOBJ.Description,
					testOBJ.Keywords,
					testOBJ.OGTitle,
					testOBJ.OGDescription,
					testOBJ.OGImage,
					testOBJ.OBJName,
					testOBJ.OBJPK,
					testOBJ.OBJName,
					testOBJ.OBJPK,
				).
				WillReturnResult(sqlmock.NewResult(1, 0))

			err := repo.UpdateSEO(context.Background(), testOBJ)
			assert.ErrorIs(t, err, rrepo.ErrNotFound)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrInternal", func(t *testing.T) {
			ErrInternal := errors.New("internal error")
			mock.ExpectExec(regexp.QuoteMeta(updateSEO)).
				WithArgs(
					testOBJ.Title,
					testOBJ.Description,
					testOBJ.Keywords,
					testOBJ.OGTitle,
					testOBJ.OGDescription,
					testOBJ.OGImage,
					testOBJ.OBJName,
					testOBJ.OBJPK,
					testOBJ.OBJName,
					testOBJ.OBJPK,
				).
				WillReturnError(ErrInternal)

			err := repo.UpdateSEO(context.Background(), testOBJ)
			assert.ErrorIs(t, err, ErrInternal)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)
}

func TestRepository_DeleteSEO(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	name, pk := "name", "pk"

	t.Run(
		"Success", func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(deleteSEO)).
				WithArgs(name, pk).
				WillReturnResult(sqlmock.NewResult(0, 1))

			err := repo.DeleteSEO(context.Background(), name, pk)
			assert.NoError(t, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrNotFound", func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(deleteSEO)).
				WithArgs(name, pk).
				WillReturnResult(sqlmock.NewResult(1, 0))

			err := repo.DeleteSEO(context.Background(), name, pk)
			assert.ErrorIs(t, err, rrepo.ErrNotFound)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"ErrInternal", func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(deleteSEO)).
				WithArgs(name, pk).
				WillReturnError(errors.New("db error"))

			err := repo.DeleteSEO(context.Background(), name, pk)
			assert.Error(t, err)
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		},
	)
}
