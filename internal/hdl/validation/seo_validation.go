package validation

import md "github.com/JMURv/seo/internal/models"

func ValidateSEO(seo *md.SEO) error {
	if seo.Title == "" {
		return ErrMissingTitle
	}

	if seo.Description == "" {
		return ErrMissingDescription
	}

	if seo.Keywords == "" {
		return ErrMissingKeywords
	}

	if seo.OGTitle == "" {
		return ErrMissingOGTitle
	}

	if seo.OGDescription == "" {
		return ErrMissingOGDescription
	}

	if seo.OGImage == "" {
		return ErrMissingOGImage
	}

	if seo.OBJName == "" {
		return ErrMissingOBJName
	}

	if seo.OBJPK == "" {
		return ErrMissingOBJPK
	}

	return nil
}
