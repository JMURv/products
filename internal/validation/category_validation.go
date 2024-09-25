package validation

import "github.com/JMURv/par-pro/products/pkg/model"

func CategoryValidation(req *model.Category) error {
	if req.Title == "" {
		return ErrMissingTitle
	}

	return nil

}
