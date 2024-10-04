package db

import (
	"context"
	"database/sql"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) GetSEO(ctx context.Context, name, pk string) (*model.SEO, error) {
	const op = "seo.GetSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.SEO{}
	err := r.conn.QueryRowContext(ctx, `
		SELECT id, title, description, keywords, og_title, og_description, og_image, obj_name, obj_pk, created_at, updated_at
		FROM seo
		WHERE obj_name = $1 AND obj_pk = $2
		`, name, pk).
		Scan(&res.ID,
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

func (r *Repository) CreateSEO(ctx context.Context, req *model.SEO) (uint64, error) {
	const op = "seo.CreateSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var idx uint64
	err := r.conn.QueryRow(`SELECT id FROM seo WHERE obj_name = $1 AND obj_pk = $2`, req.OBJName, req.OBJPK).Scan(&idx)
	if err == nil {
		return 0, repo.ErrAlreadyExists
	} else if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	var id uint64
	if err := r.conn.QueryRowContext(ctx,
		`INSERT INTO seo (
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
		 RETURNING id`,
		req.Title,
		req.Description,
		req.Keywords,
		req.OGTitle,
		req.OGDescription,
		req.OGImage,
		req.OBJName,
		req.OBJPK,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) UpdateSEO(ctx context.Context, req *model.SEO) error {
	const op = "seo.UpdateSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx,
		`UPDATE seo 
		 SET 
		 title = $1,
		 description = $2, 
		 keywords = $3, 
		 og_title = $4, 
		 og_description = $5, 
		 og_image = $6,
		 obj_name = $7, 
		 obj_pk = $8
		 WHERE obj_name = $9 AND obj_pk = $10`,
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
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.conn.ExecContext(ctx, `DELETE FROM seo WHERE obj_name = $1 AND obj_pk = $2`, name, pk)
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
