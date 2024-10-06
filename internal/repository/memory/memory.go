package memory

import (
	"github.com/JMURv/seo-svc/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	SEOData  map[uint64]*model.SEO
	PageData map[string]*model.Page
}

func New() *Repository {
	return &Repository{
		SEOData:  make(map[uint64]*model.SEO),
		PageData: make(map[string]*model.Page),
	}
}
