package memory

import (
	"github.com/JMURv/par-pro-seo/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.SEO
}

func New() *Repository {
	return &Repository{
		data: make(map[uint64]*model.SEO),
	}
}
