package validation

import "errors"

var ErrMissingTitle = errors.New("missing title")
var ErrMissingDescription = errors.New("missing description")
var ErrMissingPrice = errors.New("missing price")
var ErrMissingSrc = errors.New("missing src")
var ErrMissingUUID = errors.New("missing uuid")
