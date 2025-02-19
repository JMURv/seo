package ctrl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JMURv/seo/internal/config"
	"github.com/JMURv/seo/internal/dto"
	"github.com/JMURv/seo/internal/models"
	"github.com/JMURv/seo/internal/repo"
	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const pageKey = "page:%v"

func (c *Controller) ListPages(ctx context.Context) ([]*models.Page, error) {
	const op = "page.ListPages.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.ListPages(ctx)
	if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.Error(err),
		)
		return nil, err
	}

	return res, nil
}

func (c *Controller) GetPage(ctx context.Context, slug string) (*models.Page, error) {
	const op = "page.GetPage.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	cached := &models.Page{}
	key := fmt.Sprintf(pageKey, slug)
	if err := c.cache.GetToStruct(ctx, key, cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetPage(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			ErrNotFound.Error(),
			zap.String("op", op),
			zap.String("slug", slug),
			zap.Error(err),
		)
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.String("slug", slug),
			zap.Error(err),
		)
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		c.cache.Set(ctx, config.DefaultCacheTime, key, bytes)
	}
	return res, nil
}

func (c *Controller) CreatePage(ctx context.Context, req *models.Page) (*dto.CreatePageResponse, error) {
	const op = "page.CreatePage.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.CreatePage(ctx, req)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			ErrAlreadyExists.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		return nil, ErrAlreadyExists
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.Any("req", req),
			zap.Error(err),
		)
		return nil, err
	}

	return &dto.CreatePageResponse{
		Slug: res,
	}, nil
}

func (c *Controller) UpdatePage(ctx context.Context, slug string, req *models.Page) error {
	const op = "page.UpdatePage.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	err := c.repo.UpdatePage(ctx, slug, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			ErrNotFound.Error(),
			zap.String("op", op),
			zap.String("slug", slug), zap.Any("req", req),
			zap.Error(err),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.String("slug", slug), zap.Any("req", req),
			zap.Error(err),
		)
		return err
	}

	c.cache.Delete(ctx, fmt.Sprintf(pageKey, slug))
	return nil
}

func (c *Controller) DeletePage(ctx context.Context, slug string) error {
	const op = "page.DeletePage.ctrl"
	span, ctx := ot.StartSpanFromContext(ctx, op)
	defer span.Finish()

	err := c.repo.DeletePage(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			ErrNotFound.Error(),
			zap.String("op", op),
			zap.String("slug", slug),
			zap.Error(err),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			ErrInternal.Error(),
			zap.String("op", op),
			zap.String("slug", slug),
			zap.Error(err),
		)
		return err
	}

	c.cache.Delete(ctx, fmt.Sprintf(pageKey, slug))
	return nil
}
