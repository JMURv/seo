package db

import (
	"context"
	"database/sql"
	md "github.com/JMURv/seo/internal/models"
	"github.com/JMURv/seo/internal/repo"
	ot "github.com/opentracing/opentracing-go"
)

func (r *Repository) GetSEO(ctx context.Context, name, pk string) (*md.SEO, error) {
	const op = "seo.GetSEO.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &md.SEO{}
	err := r.conn.QueryRowContext(ctx, getSEO, name, pk).
		Scan(
			&res.ID,
			&res.Title,
			&res.Description,
			&res.Keywords,
			&res.OGTitle,
			&res.OGDescription,
			&res.OGImage,
			&res.OBJName,
			&res.OBJPK,
			&res.CreatedAt,
			&res.UpdatedAt,
		)

	if err == sql.ErrNoRows {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreateSEO(ctx context.Context, req *md.SEO) (uint64, error) {
	const op = "seo.CreateSEO.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var id uint64
	err := r.conn.QueryRowContext(
		ctx,
		createSEO,
		req.Title,
		req.Description,
		req.Keywords,
		req.OGTitle,
		req.OGDescription,
		req.OGImage,
		req.OBJName,
		req.OBJPK,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, repo.ErrAlreadyExists
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateSEO(ctx context.Context, req *md.SEO) error {
	const op = "seo.UpdateSEO.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(
		ctx,
		updateSEO,
		req.Title,
		req.Description,
		req.Keywords,
		req.OGTitle,
		req.OGDescription,
		req.OGImage,
		req.OBJName,
		req.OBJPK,
		req.OBJName,
		req.OBJPK,
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

func (r *Repository) DeleteSEO(ctx context.Context, name, pk string) error {
	const op = "seo.DeleteSEO.repo"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, deleteSEO, name, pk)
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
