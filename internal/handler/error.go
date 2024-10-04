package handler

import "errors"

var ErrInternal = errors.New("internal error")
var ErrDecodeRequest = errors.New("failed to decode request")
var ErrMethodNotAllowed = errors.New("method not allowed")
