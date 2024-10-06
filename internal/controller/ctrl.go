package ctrl

import (
	"context"
	"time"
)

type appRepo interface {
	PageRepo
	SEORepo
}

type Discovery interface {
	Register() error
	Deregister() error
}

type CacheRepo interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, t time.Duration, key string, val any) error
	Delete(ctx context.Context, key string) error
	InvalidateKeysByPattern(ctx context.Context, pattern string) error
	Close()
}

type Controller struct {
	repo  appRepo
	cache CacheRepo
}

func New(repo appRepo, cache CacheRepo) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}
