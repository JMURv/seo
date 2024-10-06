package memory

import (
	"context"
	repo "github.com/JMURv/seo-svc/internal/repository"
	"github.com/JMURv/seo-svc/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) GetSEO(ctx context.Context, name, pk string) (*model.SEO, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "seo.GetSEO.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	for _, v := range r.SEOData {
		if v.OBJName == name && v.OBJPK == pk {
			return v, nil
		}
	}

	return nil, repo.ErrNotFound
}

func (r *Repository) CreateSEO(ctx context.Context, req *model.SEO) (uint64, error) {
	const op = "seo.CreateSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if _, err := r.GetSEO(ctx, req.OBJName, req.OBJPK); err == nil {
		return 0, repo.ErrAlreadyExists
	}

	r.Lock()
	defer r.Unlock()

	req.ID = uint64(len(r.SEOData) + 1)
	r.SEOData[req.ID] = req
	return req.ID, nil
}

func (r *Repository) UpdateSEO(ctx context.Context, req *model.SEO) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "seo.UpdateSEO.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for i, v := range r.SEOData {
		if v.OBJName == req.OBJName && v.OBJPK == req.OBJPK {
			r.SEOData[i] = req
			return nil
		}
	}
	return repo.ErrNotFound
}

func (r *Repository) DeleteSEO(ctx context.Context, name, pk string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "seo.DeleteSEO.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for i, v := range r.SEOData {
		if v.OBJName == name && v.OBJPK == pk {
			delete(r.SEOData, i+1)
			return nil
		}
	}
	return repo.ErrNotFound
}
