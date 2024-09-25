package validation

import "github.com/JMURv/par-pro/products/pkg/model"

func ItemValidation(i *model.Item) error {
	if i.Title == "" {
		return ErrMissingTitle
	}

	if i.Description == "" {
		return ErrMissingDescription
	}

	if i.Price <= 0 {
		return ErrMissingPrice
	}

	if i.Src == "" {
		return ErrMissingSrc
	}
	return nil
}
