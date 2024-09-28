package ctrl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	repo "github.com/JMURv/par-pro-seo/internal/repository"
	"github.com/JMURv/par-pro-seo/pkg/consts"
	"github.com/JMURv/par-pro-seo/pkg/model"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const SEOKey = "SEO:%v:%v"

type SEORepo interface {
	GetSEO(ctx context.Context, name, pk string) (*model.SEO, error)
	CreateSEO(ctx context.Context, name, pk string, req *model.SEO) (*model.SEO, error)
	UpdateSEO(ctx context.Context, name, pk string, req *model.SEO) (*model.SEO, error)
	DeleteSEO(ctx context.Context, name, pk string) error
}

func (c *Controller) GetSEO(ctx context.Context, name, pk string) (*model.SEO, error) {
	const op = "seo.GetSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.SEO{}
	cacheKey := fmt.Sprintf(SEOKey, name, pk)
	if err := c.cache.GetToStruct(ctx, cacheKey, cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetSEO(ctx, name, pk)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to get seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, cacheKey, bytes); err != nil {
			zap.L().Debug(
				"failed to set to cache",
				zap.Error(err), zap.String("op", op),
				zap.String("name", name), zap.String("pk", pk),
			)
		}
	}
	return res, nil
}

func (c *Controller) CreateSEO(ctx context.Context, name, pk string, req *model.SEO) error {
	const op = "seo.CreateSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.CreateSEO(ctx, name, pk, req)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			"seo with this name and pk already exists",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return ErrAlreadyExists
	} else if err != nil {
		zap.L().Debug(
			"failed to create seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(SEOKey, name, pk), bytes); err != nil {
			zap.L().Debug(
				"failed to set to cache",
				zap.Error(err), zap.String("op", op),
				zap.String("name", name), zap.String("pk", pk),
			)
		}
	}
	return nil
}

func (c *Controller) UpdateSEO(ctx context.Context, name, pk string, seo *model.SEO) error {
	const op = "seo.UpdateSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.UpdateSEO(ctx, name, pk, seo)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to update seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, fmt.Sprintf(SEOKey, name, pk), bytes); err != nil {
			zap.L().Debug(
				"failed to set to cache",
				zap.Error(err), zap.String("op", op),
				zap.String("name", name), zap.String("pk", pk),
			)
		}
	}
	return nil
}

func (c *Controller) DeleteSEO(ctx context.Context, name, pk string) error {
	const op = "seo.DeleteSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	if err := c.repo.DeleteSEO(ctx, name, pk); err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to delete seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
		return err
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(SEOKey, name, pk)); err != nil {
		zap.L().Debug(
			"failed to delete from cache",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
	}

	return nil
}
