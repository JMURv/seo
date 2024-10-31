package ctrl

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

var ErrNotFoundSvc = errors.New("service not found")
var ErrCreateClient = errors.New("failed to create client")
