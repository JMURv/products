package validation

import (
	"github.com/JMURv/par-pro/products/pkg/model"
)

func ValidatePromotion(p *model.Promotion) error {
	if p.Title == "" {
		return ErrMissingTitle
	}

	if p.Description == "" {
		return ErrMissingDescription
	}

	if p.Src == "" {
		return ErrMissingSrc
	}

	return nil
}
