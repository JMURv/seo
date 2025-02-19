package dto

type CreatePageResponse struct {
	Slug string `json:"slug"`
}

type CreateSEOResponse struct {
	Name string `json:"name"`
	PK   string `json:"pk"`
}
