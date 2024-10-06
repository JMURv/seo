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

const pageKey = "page:%v"

type PageRepo interface {
	ListPages(ctx context.Context) ([]*model.Page, error)
	GetPage(ctx context.Context, slug string) (*model.Page, error)
	CreatePage(ctx context.Context, req *model.Page) (string, error)
	UpdatePage(ctx context.Context, slug string, req *model.Page) error
	DeletePage(ctx context.Context, slug string) error
}

func (c *Controller) ListPages(ctx context.Context) ([]*model.Page, error) {
	const op = "page.ListPages.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.ListPages(ctx)
	if err != nil {
		zap.L().Debug(
			"Error list pages",
			zap.Error(err), zap.String("op", op),
		)
		return nil, err
	}

	return res, nil
}

func (c *Controller) GetPage(ctx context.Context, slug string) (*model.Page, error) {
	const op = "page.GetPage.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached := &model.Page{}
	key := fmt.Sprintf(pageKey, slug)
	if err := c.cache.Get(ctx, key, cached); err == nil {
		return cached, nil
	}

	res, err := c.repo.GetPage(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"Error page not found",
			zap.Error(err), zap.String("op", op),
		)
		return nil, ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"Error get page",
			zap.Error(err), zap.String("op", op),
		)
		return nil, err
	}

	if bytes, err := json.Marshal(res); err == nil {
		if err = c.cache.Set(ctx, consts.DefaultCacheTime, key, bytes); err != nil {
			zap.L().Debug(
				"failed to set to cache",
				zap.Error(err), zap.String("op", op),
			)
		}
	}
	return res, nil
}

func (c *Controller) CreatePage(ctx context.Context, req *model.Page) (string, error) {
	const op = "page.CreatePage.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.CreatePage(ctx, req)
	if err != nil && errors.Is(err, repo.ErrAlreadyExists) {
		zap.L().Debug(
			"Error page already exists",
			zap.Error(err), zap.String("op", op),
		)
		return "", ErrAlreadyExists
	} else if err != nil {
		zap.L().Debug(
			"Error create page",
			zap.Error(err), zap.String("op", op),
		)
		return "", err
	}

	return res, nil
}

func (c *Controller) UpdatePage(ctx context.Context, slug string, req *model.Page) error {
	const op = "page.UpdatePage.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.UpdatePage(ctx, slug, req)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"Error page not found",
			zap.Error(err), zap.String("op", op),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"Error update page",
			zap.Error(err), zap.String("op", op),
		)
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(pageKey, slug)); err != nil {
		zap.L().Debug(
			"failed to delete from cache",
			zap.Error(err), zap.String("op", op),
		)
	}

	return nil
}

func (c *Controller) DeletePage(ctx context.Context, slug string) error {
	const op = "page.DeletePage.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.DeletePage(ctx, slug)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		zap.L().Debug(
			"Error page not found",
			zap.Error(err), zap.String("op", op),
		)
		return ErrNotFound
	} else if err != nil {
		zap.L().Debug(
			"Error delete page",
			zap.Error(err), zap.String("op", op),
		)
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(pageKey, slug)); err != nil {
		zap.L().Debug(
			"failed to delete from cache",
			zap.Error(err), zap.String("op", op),
		)
	}

	return nil
}
