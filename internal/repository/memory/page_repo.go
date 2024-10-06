package memory

import (
	"context"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) GetPage(ctx context.Context, slug string) (*model.Page, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "page.GetPage.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	for _, v := range r.PageData {
		if v.Slug == slug {
			return v, nil
		}
	}

	return nil, repo.ErrNotFound
}

func (r *Repository) CreatePage(ctx context.Context, req *model.Page) (string, error) {
	const op = "page.GetPage.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if _, err := r.GetPage(ctx, req.Slug); err == nil {
		return "", repo.ErrAlreadyExists
	}

	r.Lock()
	defer r.Unlock()

	r.PageData[req.Slug] = req
	return req.Slug, nil
}

func (r *Repository) UpdatePage(ctx context.Context, slug string, req *model.Page) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "page.UpdatePage.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for i, v := range r.PageData {
		if v.Slug == slug {
			r.PageData[i] = req
			return nil
		}
	}
	return repo.ErrNotFound
}

func (r *Repository) DeletePage(ctx context.Context, slug string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "DeletePage.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for _, v := range r.PageData {
		if v.Slug == slug {
			delete(r.PageData, slug)
			return nil
		}
	}
	return repo.ErrNotFound
}
