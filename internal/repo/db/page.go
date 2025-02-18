package db

import (
	"context"
	"database/sql"
	"github.com/JMURv/seo/internal/config"
	md "github.com/JMURv/seo/internal/models"
	"github.com/JMURv/seo/internal/repo"
	ot "github.com/opentracing/opentracing-go"
)

func (r *Repository) ListPages(ctx context.Context) ([]*md.Page, error) {
	const op = "pages.ListPages.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	rows, err := r.conn.QueryContext(ctx, listPage)
	if err != nil {
		return nil, err
	}

	res := make([]*md.Page, 0, config.DefaultSize)
	for rows.Next() {
		page := &md.Page{}
		if err = rows.Scan(&page.Slug, &page.Title, &page.Href, &page.CreatedAt, &page.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, page)
	}

	return res, nil
}

func (r *Repository) GetPage(ctx context.Context, slug string) (*md.Page, error) {
	const op = "pages.GetPage.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &md.Page{}
	err := r.conn.QueryRowContext(ctx, getPageBySlug, slug).
		Scan(&res.Slug, &res.Title, &res.Href, &res.CreatedAt, &res.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreatePage(ctx context.Context, req *md.Page) (string, error) {
	const op = "pages.CreatePage.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var slug string
	err := r.conn.QueryRowContext(ctx, createPage, req.Slug, req.Title, req.Href).Scan(&slug)
	if err == sql.ErrNoRows {
		return "", repo.ErrAlreadyExists
	} else if err != nil {
		return "", err
	}

	return slug, nil
}

func (r *Repository) UpdatePage(ctx context.Context, slug string, req *md.Page) error {
	const op = "pages.UpdatePage.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, updatePage, req.Title, req.Href, slug)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return repo.ErrNotFound
	}

	return nil
}

func (r *Repository) DeletePage(ctx context.Context, slug string) error {
	const op = "pages.DeletePage.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, deletePage, slug)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if aff == 0 {
		return repo.ErrNotFound
	}
	return nil
}
