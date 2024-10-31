package db

import (
	"github.com/JMURv/par-pro/products/pkg/model"
	"strconv"
	"strings"
)

func ScanChildrenCategory(req []string) ([]model.Category, error) {
	res := make([]model.Category, 0, len(req))
	for _, v := range req {
		parts := strings.Split(v, "|")
		if len(parts) != 5 {
			continue
		}

		res = append(
			res, model.Category{
				Slug:       parts[0],
				Title:      parts[1],
				Src:        parts[2],
				Alt:        parts[3],
				ParentSlug: parts[4],
			},
		)
	}
	return res, nil
}

func ScanFilters(req []string) ([]model.Filter, error) {
	res := make([]model.Filter, 0, len(req))
	for _, v := range req {
		parts := strings.Split(v, "|")
		if len(parts) != 6 {
			continue
		}

		id, err := strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			return nil, err
		}

		values := strings.Split(parts[2], ",")

		minVal, err := strconv.ParseFloat(parts[4], 64)
		if err != nil {
			return nil, err
		}

		maxVal, err := strconv.ParseFloat(parts[5], 64)
		if err != nil {
			return nil, err
		}

		res = append(
			res, model.Filter{
				ID:         id,
				Name:       parts[1],
				Values:     values,
				FilterType: parts[3],
				MinValue:   minVal,
				MaxValue:   maxVal,
			},
		)
	}
	return res, nil
}
