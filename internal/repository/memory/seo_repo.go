package memory

import (
	"context"
	repo "github.com/JMURv/par-pro-seo/internal/repository"
	"github.com/JMURv/par-pro-seo/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) GetSEO(ctx context.Context, name, pk string) (*model.SEO, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "seo.GetSEO.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	for _, v := range r.data {
		if v.OBJName == name && v.OBJPK == pk {
			return v, nil
		}
	}

	return nil, repo.ErrNotFound
}

func (r *Repository) CreateSEO(ctx context.Context, name, pk string, req *model.SEO) (*model.SEO, error) {
	const op = "seo.CreateSEO.repo"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	res, err := r.GetSEO(ctx, name, pk)
	if err == nil {
		return nil, repo.ErrAlreadyExists
	}

	r.data[uint64(len(r.data)+1)] = req
	return res, nil
}

func (r *Repository) UpdateSEO(ctx context.Context, name, pk string, req *model.SEO) (*model.SEO, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "seo.UpdateSEO.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for i, v := range r.data {
		if v.OBJName == name && v.OBJPK == pk {
			r.data[i] = req
			return req, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) DeleteSEO(ctx context.Context, name, pk string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "seo.DeleteSEO.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for i, v := range r.data {
		if v.OBJName == name && v.OBJPK == pk {
			delete(r.data, i+1)
			return nil
		}
	}
	return repo.ErrNotFound
}
