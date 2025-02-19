package db

const getSEO = `
SELECT title, description, keywords, og_title, og_description, og_image, obj_name, obj_pk, created_at, updated_at
FROM seo
WHERE obj_name = $1 AND obj_pk = $2
`

const createSEO = `
INSERT INTO seo (
	title, 
	description, 
	keywords,
	og_title,
	og_description,
	og_image,
	obj_name,
	obj_pk
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (obj_name, obj_pk) DO NOTHING
RETURNING obj_name, obj_pk
`

const updateSEO = `
UPDATE seo 
SET 
	title = $1,
	description = $2, 
	keywords = $3, 
	og_title = $4, 
	og_description = $5, 
	og_image = $6,
	obj_name = $7, 
	obj_pk = $8
WHERE obj_name = $9 AND obj_pk = $10
`

const deleteSEO = `
DELETE FROM seo 
WHERE obj_name = $1 AND obj_pk = $2
`
