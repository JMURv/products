package validation

import "errors"

var ErrMissingTitle = errors.New("missing title")
var ErrMissingDescription = errors.New("missing description")
var ErrMissingPrice = errors.New("missing price")
var ErrMissingSrc = errors.New("missing src")
var ErrMissingUUID = errors.New("missing uuid")

var ErrMissingFIO = errors.New("missing fio")
var ErrMissingTel = errors.New("missing tel")
var ErrMissingEmail = errors.New("missing email")
var ErrMissingAddress = errors.New("missing address")
