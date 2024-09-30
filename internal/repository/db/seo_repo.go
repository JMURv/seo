package db

import (
	"context"
	"errors"
	repo "github.com/JMURv/par-pro-seo/internal/repository"
	"github.com/JMURv/par-pro-seo/pkg/model"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
)

func (r *Repository) GetSEO(ctx context.Context, name, pk string) (*model.SEO, error) {
	const op = "seo.GetSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res := &model.SEO{}
	err := r.conn.Where("obj_name=? AND objpk=?", name, pk).First(res).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreateSEO(ctx context.Context, req *model.SEO) (*model.SEO, error) {
	const op = "seo.CreateSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	_, err := r.GetSEO(ctx, req.OBJName, req.OBJPK)
	if err == nil {
		return nil, repo.ErrAlreadyExists
	}

	if err := r.conn.Save(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

func (r *Repository) UpdateSEO(ctx context.Context, req *model.SEO) (*model.SEO, error) {
	const op = "seo.UpdateSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.GetSEO(ctx, req.OBJName, req.OBJPK)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		res.Title = req.Title
	}
	if req.Description != "" {
		res.Description = req.Description
	}
	if req.Keywords != "" {
		res.Keywords = req.Keywords
	}
	if req.OGTitle != "" {
		res.OGTitle = req.OGTitle
	}
	if req.OGDescription != "" {
		res.OGDescription = req.OGDescription
	}
	if req.OGImage != "" {
		res.OGImage = req.OGImage
	}
	if req.OBJName != "" {
		res.OBJName = req.OBJName
	}
	if req.OBJPK != "" {
		res.OBJPK = req.OBJPK
	}

	if err := r.conn.Save(res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repository) DeleteSEO(ctx context.Context, name, pk string) error {
	const op = "seo.DeleteSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := r.GetSEO(ctx, name, pk)
	if err != nil {
		return repo.ErrNotFound
	}

	if err := r.conn.Delete(res).Error; err != nil {
		return err
	}
	return nil
}
