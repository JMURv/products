package db

import (
	"database/sql"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
)

func createPromotionItem(tx *sql.Tx, slug string, req model.PromotionItem) error {
	if _, err := tx.Exec(promoItemCreateQ, req.Discount, slug, req.ItemID); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func CreatePromotionItems(tx *sql.Tx, slug string, req []*model.PromotionItem) error {
	for _, v := range req {
		if err := createPromotionItem(tx, slug, *v); err != nil {
			return err
		}
	}
	return nil
}

func UpdatePromotionItems(tx *sql.Tx, slug string, req []*model.PromotionItem) error {
	rows, err := tx.Query(promoItemListQ, slug)
	if err != nil {
		return err
	}
	defer rows.Close()

	existing := make(map[uuid.UUID]model.PromotionItem)
	for rows.Next() {
		var v model.PromotionItem
		if err = rows.Scan(
			&v.Discount,
			&v.PromotionSlug,
			&v.ItemID,
		); err != nil {
			return err
		}
		existing[v.ItemID] = v
	}

	if err = rows.Err(); err != nil {
		return err
	}

	for _, v := range req {
		if _, exists := existing[v.ItemID]; exists {
			if _, err = tx.Exec(
				promoItemUpdateQ,
				v.Discount,
				v.PromotionSlug,
				v.ItemID,
				slug,
			); err != nil {
				return err
			}
		} else {
			if err = createPromotionItem(tx, slug, *v); err != nil {
				return err
			}
		}
	}

	reqMap := make(map[uuid.UUID]struct{}, len(req))
	for _, v := range req {
		reqMap[v.ItemID] = struct{}{}
	}
	for id := range existing {
		if _, ok := reqMap[id]; !ok {
			if _, err = tx.Exec(promoItemDelete, id, slug); err != nil {
				return err
			}
		}
	}

	return nil
}
