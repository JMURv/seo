package ctrl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JMURv/seo/internal/config"
	md "github.com/JMURv/seo/internal/models"
	"github.com/JMURv/seo/internal/repo"
	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const SEOKey = "SEO:%v:%v"

func (c *Controller) GetSEO(ctx context.Context, name, pk string) (*md.SEO, error) {
	const op = "seo.GetSEO.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	cached := &md.SEO{}
	key := fmt.Sprintf(SEOKey, name, pk)
	if err := c.cache.GetToStruct(ctx, key, cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetSEO(ctx, name, pk)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			ErrNotFound.Error(),
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
			zap.Error(err),
		)
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
			zap.Error(err),
		)
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		c.cache.Set(ctx, config.DefaultCacheTime, key, bytes)
	}
	return res, nil
}

func (c *Controller) CreateSEO(ctx context.Context, req *md.SEO) (uint64, error) {
	const op = "seo.CreateSEO.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.CreateSEO(ctx, req)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			ErrAlreadyExists.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		return 0, ErrAlreadyExists
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		return 0, err
	}

	return res, nil
}

func (c *Controller) UpdateSEO(ctx context.Context, req *md.SEO) error {
	const op = "seo.UpdateSEO.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	err := c.repo.UpdateSEO(ctx, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			ErrNotFound.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		return err
	}

	c.cache.Delete(ctx, fmt.Sprintf(SEOKey, req.OBJName, req.OBJPK))
	return nil
}

func (c *Controller) DeleteSEO(ctx context.Context, name, pk string) error {
	const op = "seo.DeleteSEO.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	key := fmt.Sprintf(SEOKey, name, pk)
	if err := c.repo.DeleteSEO(ctx, name, pk); err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			ErrNotFound.Error(),
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
			zap.Error(err),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.String("name", name), zap.String("pk", pk),
			zap.Error(err),
		)
		return err
	}

	c.cache.Delete(ctx, key)
	return nil
}
