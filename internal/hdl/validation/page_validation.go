package validation

import md "github.com/JMURv/seo/internal/models"

func ValidatePage(req *md.Page) error {
	if req.Slug == "" {
		return ErrMissingSlug
	}

	if req.Title == "" {
		return ErrMissingTitle
	}

	if req.Href == "" {
		return ErrMissingHref
	}
	return nil
}
