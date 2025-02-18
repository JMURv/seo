package ctrl

import (
	"context"
	md "github.com/JMURv/seo/internal/models"
	"io"
	"time"
)

type AppRepo interface {
	GetSEO(ctx context.Context, name, pk string) (*md.SEO, error)
	CreateSEO(ctx context.Context, req *md.SEO) (uint64, error)
	UpdateSEO(ctx context.Context, req *md.SEO) error
	DeleteSEO(ctx context.Context, name, pk string) error

	ListPages(ctx context.Context) ([]*md.Page, error)
	GetPage(ctx context.Context, slug string) (*md.Page, error)
	CreatePage(ctx context.Context, req *md.Page) (string, error)
	UpdatePage(ctx context.Context, slug string, req *md.Page) error
	DeletePage(ctx context.Context, slug string) error
}

type AppCtrl interface {
	GetSEO(ctx context.Context, name, pk string) (*md.SEO, error)
	CreateSEO(ctx context.Context, req *md.SEO) (uint64, error)
	UpdateSEO(ctx context.Context, req *md.SEO) error
	DeleteSEO(ctx context.Context, name, pk string) error

	ListPages(ctx context.Context) ([]*md.Page, error)
	GetPage(ctx context.Context, slug string) (*md.Page, error)
	CreatePage(ctx context.Context, req *md.Page) (string, error)
	UpdatePage(ctx context.Context, slug string, req *md.Page) error
	DeletePage(ctx context.Context, slug string) error
}

type CacheService interface {
	io.Closer
	GetToStruct(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, t time.Duration, key string, val any)
	Delete(ctx context.Context, key string)
	InvalidateKeysByPattern(ctx context.Context, pattern string)
}

type Controller struct {
	repo  AppRepo
	cache CacheService
}

func New(repo AppRepo, cache CacheService) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}
