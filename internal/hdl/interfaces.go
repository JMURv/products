package hdl

import (
	"context"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
)

type Ctrl interface {
	ListFavorites(ctx context.Context, uid uuid.UUID) ([]*model.Favorite, error)
	AddToFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) (*model.Favorite, error)
	RemoveFromFavorites(ctx context.Context, uid uuid.UUID, itemID uuid.UUID) error
	ListCategoryItems(ctx context.Context, slug string, page int, size int, filters map[string]any, sort string) (*model.PaginatedItemsData, error)
	ItemAttrSearch(ctx context.Context, query string, size int, page int) (*model.PaginatedItemAttrData, error)
	ItemSearch(ctx context.Context, query string, page, size int) (*model.PaginatedItemsData, error)
	ListRelatedItems(ctx context.Context, uid uuid.UUID) ([]*model.RelatedProduct, error)
	ListItemsByLabel(ctx context.Context, label string, page int, size int) (*model.PaginatedItemsData, error)
	ListItems(ctx context.Context, page int, size int) (*model.PaginatedItemsData, error)
	GetItemByUUID(ctx context.Context, uid uuid.UUID) (*model.Item, error)
	CreateItem(ctx context.Context, i *model.Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, uid uuid.UUID, i *model.Item) error
	DeleteItem(ctx context.Context, uid uuid.UUID) error
	CategoryFiltersSearch(ctx context.Context, query string, page int, size int) (*model.PaginatedFilterData, error)
	CategorySearch(ctx context.Context, query string, page int, size int) (*model.PaginatedCategoryData, error)
	ListCategoryFilters(ctx context.Context, slug string) ([]*model.Filter, error)
	ListCategories(ctx context.Context, page int, size int) (*model.PaginatedCategoryData, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) (string, error)
	UpdateCategory(ctx context.Context, slug string, category *model.Category) error
	DeleteCategory(ctx context.Context, slug string) error
	PromotionSearch(ctx context.Context, query string, page int, size int) (*model.PaginatedPromosData, error)
	ListPromotionItems(ctx context.Context, slug string, page int, size int) (*model.PaginatedPromoItemsData, error)
	ListPromotions(ctx context.Context, page int, size int) (*model.PaginatedPromosData, error)
	GetPromotion(ctx context.Context, slug string) (*model.Promotion, error)
	CreatePromotion(ctx context.Context, p *model.Promotion) (string, error)
	UpdatePromotion(ctx context.Context, slug string, p *model.Promotion) error
	DeletePromotion(ctx context.Context, slug string) error
}
