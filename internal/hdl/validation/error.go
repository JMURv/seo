package validation

import "errors"

var ErrMissingSlug = errors.New("missing slug")
var ErrMissingTitle = errors.New("missing title")
var ErrMissingDescription = errors.New("missing description")
var ErrMissingKeywords = errors.New("missing keywords")
var ErrMissingOGTitle = errors.New("missing og title")
var ErrMissingOGDescription = errors.New("missing og description")
var ErrMissingOGImage = errors.New("missing og image")
var ErrMissingOBJName = errors.New("missing related obj name")
var ErrMissingOBJPK = errors.New("missing related obj pk")

var ErrMissingHref = errors.New("missing href")
