package validation

import "github.com/JMURv/seo-svc/pkg/model"

func ValidatePage(req *model.Page) error {
	if req.Slug == "" {
		return ErrMissingDescription
	}

	if req.Title == "" {
		return ErrMissingTitle
	}

	if req.Href == "" {
		return ErrMissingHref
	}
	return nil
}
