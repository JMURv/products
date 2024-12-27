package db

import (
	"database/sql"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
)

func CreateItemMedia(tx *sql.Tx, uid uuid.UUID, media []model.ItemMedia) error {
	for _, v := range media {
		if _, err := tx.Exec(itemMediaCreateQ, uid, v.Src, v.Alt); err != nil {
			return err
		}
	}
	return nil
}

func CreateItemAttrs(tx *sql.Tx, uid uuid.UUID, attrs []model.ItemAttribute) error {
	for _, v := range attrs {
		if _, err := tx.Exec(itemAttrCreateQ, uid, v.Name, v.Value); err != nil {
			return err
		}
	}
	return nil
}

func CreateItemCategories(tx *sql.Tx, uid uuid.UUID, categories []model.Category) error {
	for _, v := range categories {
		if _, err := tx.Exec(itemCategoryCreateQ, uid, v.Slug); err != nil {
			return err
		}
	}
	return nil
}

func UpdateItemMedia(tx *sql.Tx, uid uuid.UUID, req []model.ItemMedia) error {
	rows, err := tx.Query(itemMediaList, uid)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[string]model.ItemMedia)
	for rows.Next() {
		var media model.ItemMedia
		if err = rows.Scan(&media.ID, &media.Src, &media.Alt); err != nil {
			return err
		}
		existing[media.Src] = media
	}

	if err = rows.Err(); err != nil {
		return err
	}

	reqMap := make(map[string]struct{}, len(req))
	for _, v := range req {
		reqMap[v.Src] = struct{}{}
		if _, ok := existing[v.Src]; !ok {
			if _, err = tx.Exec(itemMediaCreateQ, uid, v.Src, v.Alt); err != nil {
				return err
			}
		}
	}

	for src, media := range existing {
		if _, ok := reqMap[src]; !ok {
			if _, err = tx.Exec(itemMediaDeleteQ, media.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateItemAttributes(tx *sql.Tx, uid uuid.UUID, req []model.ItemAttribute) error {
	rows, err := tx.Query(itemAttrList, uid)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[string]model.ItemAttribute)
	for rows.Next() {
		var attr model.ItemAttribute
		if err = rows.Scan(&attr.ID, &attr.Name, &attr.Value); err != nil {
			return err
		}
		existing[attr.Name] = attr
	}

	if err = rows.Err(); err != nil {
		return err
	}

	for _, v := range req {
		if attr, exists := existing[v.Name]; exists {
			if attr.Name != v.Name || attr.Value != v.Value {
				if _, err = tx.Exec(itemAttrUpdateQ, v.Name, v.Value, attr.ID); err != nil {
					return err
				}
			}
		} else {
			if _, err = tx.Exec(itemAttrCreateQ, uid, v.Name, v.Value); err != nil {
				return err
			}
		}
	}

	reqMap := make(map[string]struct{}, len(req))
	for _, v := range req {
		reqMap[v.Name] = struct{}{}
	}
	for key := range existing {
		if _, exists := reqMap[key]; !exists {
			if _, err = tx.Exec(itemAttrDeleteQ, key, uid); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateItemCategories(tx *sql.Tx, uid uuid.UUID, req []model.Category) error {
	rows, err := tx.Query(itemCategoriesListQ, uid)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[string]struct{})
	for rows.Next() {
		var ic model.ItemCategory
		if err = rows.Scan(&ic.ItemID, &ic.CategorySlug); err != nil {
			return err
		}
		existing[ic.CategorySlug] = struct{}{}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	reqMap := make(map[string]struct{}, len(req))
	for _, v := range req {
		reqMap[v.Slug] = struct{}{}
	}

	for slug := range existing {
		if _, exists := reqMap[slug]; !exists {
			if _, err = tx.Exec(itemCategoryDeleteQ, uid, slug); err != nil {
				return err
			}
		}
	}

	for slug := range reqMap {
		if _, exists := existing[slug]; !exists {
			if _, err = tx.Exec(itemCategoryCreateQ, uid, slug); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateItemRelatedProducts(tx *sql.Tx, uid uuid.UUID, req []model.RelatedProduct) error {
	rows, err := tx.Query(itemRelatedProductList, uid)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[uuid.UUID]struct{})
	for rows.Next() {
		var related model.RelatedProduct
		if err = rows.Scan(&related.ItemID, &related.RelatedItemID); err != nil {
			return err
		}
		existing[related.RelatedItemID] = struct{}{}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	reqMap := make(map[uuid.UUID]struct{}, len(req))
	for _, v := range req {
		reqMap[v.RelatedItemID] = struct{}{}
	}

	for id := range existing {
		if _, exists := reqMap[id]; !exists {
			if _, err = tx.Exec(itemRelatedProductDeleteQ, uid, id); err != nil {
				return err
			}
		}
	}

	for id := range reqMap {
		if _, exists := existing[id]; !exists {
			if _, err = tx.Exec(itemRelatedProductCreateQ, uid, id); err != nil {
				return err
			}
		}
	}

	return nil
}

func UpdateItemVariants(tx *sql.Tx, uid uuid.UUID, req []model.Item) error {
	reqMap := make(map[uuid.UUID]struct{}, len(req))
	for _, v := range req {
		reqMap[v.ID] = struct{}{}
	}

	existingVars, err := getVariantsByParentID(tx, uid)
	if err != nil {
		return err
	}

	existing := make(map[uuid.UUID]struct{})
	for _, variant := range existingVars {
		existing[variant.ID] = struct{}{}
	}

	for varID := range existing {
		if _, ok := reqMap[varID]; !ok {
			if _, err = tx.Exec(itemDeleteByParentQ, uid, varID); err != nil {
				return err
			}
		}
	}

	for _, v := range req {
		if _, exists := existing[v.ID]; exists {
			if _, err = tx.Exec(
				itemUpdateQ,
				v.Title,
				v.Article,
				v.Description,
				v.Price,
				v.Src,
				v.Alt,
				v.InStock,
				v.IsHit,
				v.IsRec,
				uid,
				v.ID,
			); err != nil {
				return err
			}

			if err = UpdateItemMedia(tx, v.ID, v.Media); err != nil {
				return err
			}

			if err = UpdateItemAttributes(tx, v.ID, v.Attributes); err != nil {
				return err
			}
		} else {
			var newVarID uuid.UUID
			if err = tx.QueryRow(
				itemCreateQ,
				v.Title,
				v.Article,
				v.Description,
				v.Price,
				v.Src,
				v.Alt,
				v.InStock,
				v.IsHit,
				v.IsRec,
				uid,
			).Scan(&newVarID); err != nil {
				return err
			}

			if err = CreateItemMedia(tx, newVarID, v.Media); err != nil {
				return err
			}

			if err = CreateItemAttrs(tx, newVarID, v.Attributes); err != nil {
				return err
			}

		}
	}

	return nil
}

func getVariantsByParentID(tx *sql.Tx, parentID uuid.UUID) ([]model.Item, error) {
	rows, err := tx.Query(
		"SELECT id, title, article, description, price, src, alt, in_stock, is_hit, is_rec FROM item WHERE parent_id = $1",
		parentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []model.Item
	for rows.Next() {
		var v model.Item
		if err = rows.Scan(
			&v.ID,
			&v.Title,
			&v.Article,
			&v.Description,
			&v.Price,
			&v.Src,
			&v.Alt,
			&v.InStock,
			&v.IsHit,
			&v.IsRec,
		); err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}
