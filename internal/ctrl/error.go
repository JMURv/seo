package ctrl

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

var ErrCreateClient = errors.New("failed to create client")
var ErrInternal = errors.New("internal error")
