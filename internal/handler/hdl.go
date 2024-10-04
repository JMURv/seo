package handler

import (
	"context"
	"github.com/JMURv/seo-svc/pkg/model"
)

type SEOCtrl interface {
	GetSEO(ctx context.Context, name, pk string) (*model.SEO, error)
	CreateSEO(ctx context.Context, req *model.SEO) (uint64, error)
	UpdateSEO(ctx context.Context, req *model.SEO) error
	DeleteSEO(ctx context.Context, name, pk string) error
}
