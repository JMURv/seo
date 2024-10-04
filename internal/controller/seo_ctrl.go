package ctrl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/consts"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const SEOKey = "SEO:%v:%v"

type SEORepo interface {
	GetSEO(ctx context.Context, name, pk string) (*model.SEO, error)
	CreateSEO(ctx context.Context, req *model.SEO) (uint64, error)
	UpdateSEO(ctx context.Context, req *model.SEO) error
	DeleteSEO(ctx context.Context, name, pk string) error
}

func (c *Controller) GetSEO(ctx context.Context, name, pk string) (*model.SEO, error) {
	const op = "seo.GetSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.SEO{}
	key := fmt.Sprintf(SEOKey, name, pk)
	if err := c.cache.Get(ctx, key, cached); err == nil {
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
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, key, bytes); err != nil {
			zap.L().Debug(
				"failed to set to cache",
				zap.Error(err), zap.String("op", op),
				zap.String("name", name), zap.String("pk", pk),
			)
		}
	}
	return res, nil
}

func (c *Controller) CreateSEO(ctx context.Context, req *model.SEO) (uint64, error) {
	const op = "seo.CreateSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.CreateSEO(ctx, req)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			"seo with this name and pk already exists",
			zap.Error(err), zap.String("op", op),
			zap.String("name", req.OBJName), zap.String("pk", req.OBJPK),
		)
		return 0, ErrAlreadyExists
	} else if err != nil {
		zap.L().Debug(
			"failed to create seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", req.OBJName), zap.String("pk", req.OBJPK),
		)
		return 0, err
	}

	return res, nil
}

func (c *Controller) UpdateSEO(ctx context.Context, req *model.SEO) error {
	const op = "seo.UpdateSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.UpdateSEO(ctx, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"failed to find seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", req.OBJName), zap.String("pk", req.OBJPK),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"failed to update seo",
			zap.Error(err), zap.String("op", op),
			zap.String("name", req.OBJName), zap.String("pk", req.OBJPK),
		)
		return err
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(SEOKey, req.OBJName, req.OBJPK)); err != nil {
		zap.L().Debug(
			"failed to delete from cache",
			zap.Error(err), zap.String("op", op),
			zap.String("name", req.OBJName), zap.String("pk", req.OBJPK),
		)
	}
	return nil
}

func (c *Controller) DeleteSEO(ctx context.Context, name, pk string) error {
	const op = "seo.DeleteSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	key := fmt.Sprintf(SEOKey, name, pk)
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

	if err := c.cache.Delete(ctx, key); err != nil {
		zap.L().Debug(
			"failed to delete from cache",
			zap.Error(err), zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
		)
	}

	return nil
}
