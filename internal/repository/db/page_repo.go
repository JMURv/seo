package db

import (
	"context"
	"database/sql"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) ListPages(ctx context.Context) ([]*model.Page, error) {
	const op = "pages.ListPages.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	rows, err := r.conn.QueryContext(ctx, `SELECT slug, title, href, created_at, updated_at FROM page`)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Page, 0)
	for rows.Next() {
		page := &model.Page{}
		if err := rows.Scan(&page.Slug, &page.Title, &page.Href, &page.CreatedAt, &page.UpdatedAt); err != nil {
			return nil, err
		}
		res = append(res, page)
	}

	return res, nil
}

func (r *Repository) GetPage(ctx context.Context, slug string) (*model.Page, error) {
	const op = "pages.GetPage.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.Page{}
	err := r.conn.QueryRowContext(ctx,
		`SELECT slug, title, href, created_at, updated_at FROM page WHERE slug = $1`, slug).
		Scan(&res.Slug, &res.Title, &res.Href, &res.CreatedAt, &res.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreatePage(ctx context.Context, req *model.Page) (string, error) {
	const op = "pages.CreatePage.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var exists string
	err := r.conn.QueryRow(`SELECT id FROM page WHERE slug = $1`, req.Slug).Scan(&exists)
	if err == nil {
		return "", repo.ErrAlreadyExists
	} else if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	var slug string
	if err := r.conn.QueryRowContext(ctx,
		`INSERT INTO page (slug, title, href) 
		 VALUES ($1, $2, $3)
		 RETURNING slug`,
		req.Slug,
		req.Title,
		req.Href,
	).Scan(&slug); err != nil {
		return "", err
	}

	return slug, nil
}

func (r *Repository) UpdatePage(ctx context.Context, slug string, req *model.Page) error {
	const op = "pages.UpdatePage.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, `UPDATE page SET title = $1, href = $2 WHERE slug = $3`,
		req.Title, req.Href, slug,
	)
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
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, `DELETE FROM page WHERE slug = $1`, slug)
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
