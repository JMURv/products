package validation

import (
	"github.com/JMURv/par-pro/products/pkg/model"
)

func Order(o *model.Order) error {
	if o.FIO == "" {
		return ErrMissingFIO
	}

	if o.Tel == "" {
		return ErrMissingTel
	}

	if o.Email == "" {
		return ErrMissingEmail
	}

	if o.Address == "" {
		return ErrMissingAddress
	}

	//if o.Delivery == "" {
	//	return ErrMissingDeliveryType
	//}
	//
	//if o.PaymentMethod == "" {
	//	return ErrMissingPaymentMethod
	//}

	return nil
}
