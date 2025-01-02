package db

import (
	"fmt"
	"github.com/JMURv/par-pro/products/pkg/model"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

func ScanOrderItems(orderItems []string) ([]model.OrderItem, error) {
	items := make([]model.OrderItem, 0, len(orderItems))

	for i := 0; i < len(orderItems); i++ {
		parts := strings.Split(orderItems[i], "|")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid order item format")
		}

		id, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}

		itemID, err := uuid.Parse(parts[1])
		if err != nil {
			return nil, err
		}

		quantity, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		items = append(
			items, model.OrderItem{
				ID:       id,
				ItemID:   itemID,
				Quantity: quantity,
			},
		)
	}
	return items, nil
}
