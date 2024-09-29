package memory

import (
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"sync"
)

type Repository struct {
	sync.RWMutex
	itemsData map[uuid.UUID]*model.Item
}

func New() *Repository {
	return &Repository{
		itemsData: make(map[uuid.UUID]*model.Item),
	}
}
