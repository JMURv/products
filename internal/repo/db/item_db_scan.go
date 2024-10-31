package db

import (
	"github.com/JMURv/par-pro/products/pkg/model"
	"strings"
)

func ScanMedia(req []string) ([]model.ItemMedia, error) {
	res := make([]model.ItemMedia, 0, len(req))
	for _, v := range req {
		parts := strings.Split(v, "|")
		if len(parts) != 2 {
			continue
		}

		res = append(
			res, model.ItemMedia{
				Src: parts[0],
				Alt: parts[1],
			},
		)
	}
	return res, nil
}

func ScanAttrs(req []string) ([]model.ItemAttribute, error) {
	res := make([]model.ItemAttribute, 0, len(req))
	for _, v := range req {
		parts := strings.Split(v, "|")
		if len(parts) != 2 {
			continue
		}

		res = append(
			res, model.ItemAttribute{
				Name:  parts[0],
				Value: parts[1],
			},
		)
	}
	return res, nil
}

func ScanItemCategories(req []string) ([]model.Category, error) {
	res := make([]model.Category, 0, len(req))
	for _, v := range req {
		parts := strings.Split(v, "|")
		if len(parts) != 2 {
			continue
		}

		res = append(
			res, model.Category{
				Title: parts[0],
				Slug:  parts[1],
			},
		)
	}
	return res, nil
}
