package ctrl

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")
var ErrInternalError = errors.New("internal error")
var ErrParseUUID = errors.New("failed to parse uuid")
var ErrDecodeRequest = errors.New("failed to decode request")

var ErrNotFoundSvc = errors.New("service not found")
var ErrCreateClient = errors.New("failed to create client")
