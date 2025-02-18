package db

//func TestRepository_ListPages(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := Repository{conn: db}
//	expectedPages := []*md.Page{
//		{
//			Slug:      "page-slug-1",
//			Title:     "Page Title 1",
//			Href:      "/page-1",
//			CreatedAt: time.Now(),
//			UpdatedAt: time.Now(),
//		},
//		{
//			Slug:      "page-slug-2",
//			Title:     "Page Title 2",
//			Href:      "/page-2",
//			CreatedAt: time.Now(),
//			UpdatedAt: time.Now(),
//		},
//	}
//	t.Run(
//		"Success", func(t *testing.T) {
//			rows := sqlmock.NewRows([]string{"slug", "title", "href", "created_at", "updated_at"}).
//				AddRow(
//					expectedPages[0].Slug,
//					expectedPages[0].Title,
//					expectedPages[0].Href,
//					expectedPages[0].CreatedAt,
//					expectedPages[0].UpdatedAt,
//				).
//				AddRow(
//					expectedPages[1].Slug,
//					expectedPages[1].Title,
//					expectedPages[1].Href,
//					expectedPages[1].CreatedAt,
//					expectedPages[1].UpdatedAt,
//				)
//
//			mock.ExpectQuery(regexp.QuoteMeta(listPage)).WillReturnRows(rows)
//
//			res, err := repo.ListPages(context.Background())
//			assert.NoError(t, err)
//			assert.Equal(t, expectedPages, res)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"QueryError", func(t *testing.T) {
//			mock.ExpectQuery(regexp.QuoteMeta(listPage)).
//				WillReturnError(errors.New("query failed"))
//
//			res, err := repo.ListPages(context.Background())
//			assert.Error(t, err)
//			assert.Nil(t, res)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ScanError", func(t *testing.T) {
//			rows := sqlmock.NewRows([]string{"slug", "title", "href", "created_at", "updated_at"}).
//				AddRow("invalid-slug", "Page Title", "/page", "invalid-created-at", time.Now())
//
//			mock.ExpectQuery(regexp.QuoteMeta(listPage)).
//				WillReturnRows(rows)
//
//			res, err := repo.ListPages(context.Background())
//			assert.Error(t, err)
//			assert.Nil(t, res)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//}
//
//func TestRepository_GetPage(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := Repository{conn: db}
//	ctx := context.Background()
//	q := regexp.QuoteMeta(getPageBySlug)
//
//	slug := "slug"
//	testOBJ := md.Page{
//		Slug:  "slug",
//		Title: "title",
//		Href:  "href",
//	}
//
//	t.Run(
//		"Success case", func(t *testing.T) {
//			mock.ExpectQuery(q).
//				WithArgs(slug).
//				WillReturnRows(
//					sqlmock.NewRows(
//						[]string{
//							"slug",
//							"title",
//							"href",
//							"created_at",
//							"updated_at",
//						},
//					).
//						AddRow(
//							testOBJ.Slug,
//							testOBJ.Title,
//							testOBJ.Href,
//							testOBJ.CreatedAt,
//							testOBJ.UpdatedAt,
//						),
//				)
//
//			result, err := repo.GetPage(ctx, slug)
//
//			assert.NoError(t, err)
//			assert.NotNil(t, result)
//			assert.Equal(t, testOBJ.Slug, result.Slug)
//			assert.Equal(t, testOBJ.Title, result.Title)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrNotFound", func(t *testing.T) {
//			mock.ExpectQuery(q).
//				WithArgs(slug).
//				WillReturnError(sql.ErrNoRows)
//
//			result, err := repo.GetPage(ctx, slug)
//			assert.Nil(t, result)
//			assert.Equal(t, rrepo.ErrNotFound, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"Unexpected error case", func(t *testing.T) {
//			notExpectedError := errors.New("not expected error")
//
//			mock.ExpectQuery(q).
//				WithArgs(slug).
//				WillReturnError(notExpectedError)
//
//			result, err := repo.GetPage(ctx, slug)
//			assert.Nil(t, result)
//			assert.Equal(t, notExpectedError, err)
//			assert.NotEqual(t, rrepo.ErrNotFound, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//}
//
//func TestRepository_CreatePage(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := Repository{conn: db}
//
//	ctx := context.Background()
//	slug := "slug"
//	testOBJ := &md.Page{
//		Slug:  slug,
//		Title: "title",
//		Href:  "href",
//	}
//
//	t.Run(
//		"Success case", func(t *testing.T) {
//
//			mock.ExpectQuery(regexp.QuoteMeta(createPage)).
//				WillReturnRows(
//					sqlmock.NewRows([]string{"slug"}).
//						AddRow(testOBJ.Slug),
//				)
//
//			_, err := repo.CreatePage(ctx, testOBJ)
//			assert.NoError(t, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrAlreadyExists", func(t *testing.T) {
//			_, err := repo.CreatePage(ctx, testOBJ)
//			assert.Equal(t, rrepo.ErrAlreadyExists, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrInternal", func(t *testing.T) {
//			internalErr := errors.New("internal error")
//			mock.ExpectQuery(regexp.QuoteMeta(createPage)).
//				WillReturnError(internalErr)
//
//			_, err := repo.CreatePage(ctx, testOBJ)
//			assert.Equal(t, internalErr, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//}
//
//func TestRepository_UpdatePage(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := Repository{conn: db}
//	slug := "slug"
//	testOBJ := &md.Page{
//		Slug:  "slug",
//		Title: "title",
//		Href:  "href",
//	}
//
//	const selectQ = `SELECT id FROM page WHERE slug = $1`
//	const insertQ = `INSERT INTO page (slug, title, href) VALUES ($1, $2, $3) RETURNING slug`
//	const updateQ = `UPDATE page SET title = $1, href = $2 WHERE slug = $3`
//
//	mock.ExpectQuery(regexp.QuoteMeta(selectQ)).
//		WithArgs(slug).
//		WillReturnError(sql.ErrNoRows)
//
//	mock.ExpectQuery(regexp.QuoteMeta(insertQ)).
//		WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(slug))
//
//	_, err = repo.CreatePage(context.Background(), testOBJ)
//	require.NoError(t, err)
//
//	t.Run(
//		"Success", func(t *testing.T) {
//			mock.ExpectExec(regexp.QuoteMeta(updateQ)).
//				WithArgs(testOBJ.Title, testOBJ.Href, testOBJ.Slug).
//				WillReturnResult(sqlmock.NewResult(1, 1))
//
//			err := repo.UpdatePage(context.Background(), slug, testOBJ)
//			assert.NoError(t, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrNotFound", func(t *testing.T) {
//			mock.ExpectExec(regexp.QuoteMeta(updateQ)).
//				WithArgs(testOBJ.Title, testOBJ.Href, testOBJ.Slug).
//				WillReturnResult(sqlmock.NewResult(1, 0))
//
//			err := repo.UpdatePage(context.Background(), slug, testOBJ)
//			assert.ErrorIs(t, err, rrepo.ErrNotFound)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrInternal", func(t *testing.T) {
//			ErrInternal := errors.New("internal error")
//			mock.ExpectExec(regexp.QuoteMeta(updateQ)).
//				WithArgs(testOBJ.Title, testOBJ.Href, testOBJ.Slug).
//				WillReturnError(ErrInternal)
//
//			err := repo.UpdatePage(context.Background(), slug, testOBJ)
//			assert.ErrorIs(t, err, ErrInternal)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//}
//
//func TestRepository_DeletePage(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	require.NoError(t, err)
//	defer db.Close()
//
//	repo := Repository{conn: db}
//	slug := "slug"
//	const deleteQ = `DELETE FROM page WHERE slug = $1`
//
//	t.Run(
//		"Success", func(t *testing.T) {
//			mock.ExpectExec(regexp.QuoteMeta(deleteQ)).
//				WithArgs(slug).
//				WillReturnResult(sqlmock.NewResult(0, 1))
//
//			err := repo.DeletePage(context.Background(), slug)
//			assert.NoError(t, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrNotFound", func(t *testing.T) {
//			mock.ExpectExec(regexp.QuoteMeta(deleteQ)).
//				WithArgs(slug).
//				WillReturnResult(sqlmock.NewResult(1, 0))
//
//			err := repo.DeletePage(context.Background(), slug)
//			assert.ErrorIs(t, err, rrepo.ErrNotFound)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//
//	t.Run(
//		"ErrInternal", func(t *testing.T) {
//			mock.ExpectExec(regexp.QuoteMeta(deleteQ)).
//				WithArgs(slug).
//				WillReturnError(errors.New("db error"))
//
//			err := repo.DeletePage(context.Background(), slug)
//			assert.Error(t, err)
//			err = mock.ExpectationsWereMet()
//			assert.NoError(t, err)
//		},
//	)
//}
