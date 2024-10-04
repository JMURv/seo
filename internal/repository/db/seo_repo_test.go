package db

import (
	"context"
	"database/sql"
	"errors"
	rrepo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
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
	q := regexp.QuoteMeta(
		`SELECT id, title, description, keywords, og_title, og_description, og_image, obj_name, obj_pk, created_at, updated_at
		FROM seo
		WHERE obj_name = $1 AND obj_pk = $2`)

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

	t.Run("Success case", func(t *testing.T) {
		mock.ExpectQuery(q).
			WithArgs(name, pk).
			WillReturnRows(sqlmock.NewRows([]string{
				"id",
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
			}).
				AddRow(
					testOBJ.ID,
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
				))

		result, err := repo.GetSEO(ctx, name, pk)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testOBJ.ID, result.ID)
		assert.Equal(t, testOBJ.Title, result.Title)
		assert.Equal(t, testOBJ.Description, result.Description)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mock.ExpectQuery(q).
			WithArgs(name, pk).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetSEO(ctx, name, pk)
		assert.Nil(t, result)
		assert.Equal(t, rrepo.ErrNotFound, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("Unexpected error case", func(t *testing.T) {
		notExpectedError := errors.New("not expected error")

		mock.ExpectQuery(q).
			WithArgs(name, pk).
			WillReturnError(notExpectedError)

		result, err := repo.GetSEO(ctx, name, pk)
		assert.Nil(t, result)
		assert.Equal(t, notExpectedError, err)
		assert.NotEqual(t, rrepo.ErrNotFound, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestRepository_CreateSEO(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}

	const selectQ = `SELECT id FROM seo WHERE obj_name = $1 AND obj_pk = $2`
	const q = `INSERT INTO seo (
			 title, 
			 description, 
			 keywords,
			 og_title,
			 og_description,
			 og_image,
			 obj_name,
			 obj_pk
		 ) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`

	ctx := context.Background()
	name, pk := "name", "pk"
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

	t.Run("Success case", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectQ)).
			WithArgs(testOBJ.OBJName, testOBJ.OBJPK).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(
			regexp.QuoteMeta(q)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).
					AddRow(testOBJ.ID))

		_, err := repo.CreateSEO(ctx, testOBJ)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectQ)).
			WithArgs(testOBJ.OBJName, testOBJ.OBJPK).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(1))

		_, err := repo.CreateSEO(ctx, testOBJ)
		assert.Equal(t, rrepo.ErrAlreadyExists, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ErrSelect", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectQ)).
			WithArgs(testOBJ.OBJName, testOBJ.OBJPK).
			WillReturnError(errors.New("select error"))

		_, err := repo.CreateSEO(ctx, testOBJ)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		internalErr := errors.New("internal error")
		mock.ExpectQuery(regexp.QuoteMeta(selectQ)).
			WithArgs(testOBJ.OBJName, testOBJ.OBJPK).
			WillReturnError(sql.ErrNoRows)

		mock.ExpectQuery(regexp.QuoteMeta(q)).
			WillReturnError(internalErr)

		_, err := repo.CreateSEO(ctx, testOBJ)
		assert.Equal(t, internalErr, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
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

	const selectQ = `SELECT id FROM seo WHERE obj_name = $1 AND obj_pk = $2`
	const insertQ = `INSERT INTO seo (
			 title, 
			 description, 
			 keywords,
			 og_title,
			 og_description,
			 og_image,
			 obj_name,
			 obj_pk
		 ) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`
	const updateQ = `UPDATE seo 
		 SET 
		 title = $1,
		 description = $2, 
		 keywords = $3, 
		 og_title = $4, 
		 og_description = $5, 
		 og_image = $6,
		 obj_name = $7, 
		 obj_pk = $8
		 WHERE obj_name = $9 AND obj_pk = $10`

	mock.ExpectQuery(regexp.QuoteMeta(selectQ)).
		WithArgs(testOBJ.OBJName, testOBJ.OBJPK).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(regexp.QuoteMeta(insertQ)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testOBJ.ID))

	_, err = repo.CreateSEO(context.Background(), testOBJ)
	require.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(updateQ)).
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
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(updateQ)).
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
	})

	t.Run("ErrInternal", func(t *testing.T) {
		ErrInternal := errors.New("internal error")
		mock.ExpectExec(regexp.QuoteMeta(updateQ)).
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
	})
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := Repository{conn: db}
	name, pk := "name", "pk"
	const deleteQ = `DELETE FROM seo WHERE obj_name = $1 AND obj_pk = $2`

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(deleteQ)).
			WithArgs(name, pk).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteSEO(context.Background(), name, pk)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(deleteQ)).
			WithArgs(name, pk).
			WillReturnResult(sqlmock.NewResult(1, 0))

		err := repo.DeleteSEO(context.Background(), name, pk)
		assert.ErrorIs(t, err, rrepo.ErrNotFound)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("ErrInternal", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(deleteQ)).
			WithArgs(name, pk).
			WillReturnError(errors.New("db error"))

		err := repo.DeleteSEO(context.Background(), name, pk)
		assert.Error(t, err)
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
