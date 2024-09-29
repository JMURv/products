package utils

import (
	"fmt"
	"gorm.io/gorm"
)

func FilterOrders(query *gorm.DB, filters map[string]any) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "min_price":
			query = query.Where("total_amount >= ?", value)
			continue
		case "max_price":
			query = query.Where("total_amount <= ?", value)
			continue
		}

		switch v := value.(type) {
		case []string:
			query = query.Where(key+" IN (?)", v)
		}

	}
	return query
}

func FilterItems(query *gorm.DB, filters map[string]any) *gorm.DB {
	counter := 0
	for key, value := range filters {
		switch key {
		case "min_price":
			query = query.Where("price >= ?", value)
			continue
		case "max_price":
			query = query.Where("price <= ?", value)
			continue
		}

		switch v := value.(type) {

		case map[string]any:
			counter++
			attrAlias := fmt.Sprintf("ia%d", counter)

			min, minOK := v["min"].(string)
			max, maxOK := v["max"].(string)

			query = query.Joins(fmt.Sprintf("LEFT JOIN item_attributes %s ON %s.item_id = items.id", attrAlias, attrAlias))
			if minOK {
				query = query.Where(fmt.Sprintf("(%s.name = ? AND %s.value >= ?)", attrAlias, attrAlias), key, min)
			}
			if maxOK {
				query = query.Where(fmt.Sprintf("(%s.name = ? AND %s.value <= ?)", attrAlias, attrAlias), key, max)
			}

		case []string:
			counter++
			attrAlias := fmt.Sprintf("ia%d", counter)
			query = query.
				Joins(fmt.Sprintf("LEFT JOIN item_attributes %s ON %s.item_id = items.id", attrAlias, attrAlias)).
				Where(fmt.Sprintf("(%s.name = ? AND %s.value IN (?))", attrAlias, attrAlias), key, v)
		}

	}
	return query
}
