package db

const listPage = `
SELECT slug, title, href, created_at, updated_at 
FROM page
`

const getPageBySlug = `
SELECT slug, title, href, created_at, updated_at 
FROM page
WHERE slug = $1
`

const createPage = `
INSERT INTO page (slug, title, href) 
VALUES ($1, $2, $3)
ON CONFLICT (slug) DO NOTHING 
RETURNING slug
`

const updatePage = `
UPDATE page 
SET title = $1, href = $2 
WHERE slug = $3
`

const deletePage = `
DELETE FROM page 
WHERE slug = $1
`
