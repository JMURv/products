package db

const filterCountQ = "SELECT COUNT(*) FROM filter WHERE name ILIKE ?;"
const filterSearchQ = `
	SELECT 
	    id, 
	    name, 
	    values, 
	    filter_type, 
	    min_value, 
	    max_value, 
	    category_slug 
	FROM filter 
	WHERE name ILIKE $1
	OFFSET $2 LIMIT $3;
`

const filterListQ = `
	SELECT 
	    id,
	    name,
	    values, 
	    filter_type, 
	    min_value, 
	    max_value, 
	    category_slug 
	FROM filter
	WHERE category_slug = $1
`

const filterUpdateQ = `
	UPDATE filter
	SET 
	    name = $1, 
	    values = $2,
	    filter_type = $3,
	    min_value = $4,
	    max_value = $5,
	    category_slug = $6
	WHERE id = $6
`

const categoryCountQ = "SELECT COUNT(*) FROM category;"
const categorySearchCountQ = `SELECT COUNT(*) FROM category WHERE title ILIKE ? OR slug ILIKE ?`

const categorySearchQ = `
	SELECT slug, title, src, alt, parent_slug 
	FROM category
	WHERE title ILIKE $1 OR slug ILIKE $2
	OFFSET $3 LIMIT $4;
`

const categoryListQ = `
	SELECT 
	    p.slug, p.title, p.src, p.alt, p.parent_slug,
	    ARRAY_AGG(
             child.slug || '|' || child.title || '|' || child.src || '|' || child.alt || '|' || child.parent_slug
		) AS children
	FROM category AS p
	JOIN category AS child 
	ON child.parent_slug = p.slug
	GROUP BY p.slug, p.title
	OFFSET $1 LIMIT $2;
`

const categoryGetQ = `
	SELECT c.slug, c.title, c.src, c.alt, c.parent_slug,
	ARRAY_AGG(
		child.slug || '|' || child.title || '|' || child.src || '|' || child.alt || '|' || child.parent_slug
	) AS children,
	ARRAY_AGG(
		f.id || '|' || f.name || '|' || f.values || '|' || f.filter_type || '|' || f.min_value || '|' || f.max_value
	) AS filters
	FROM category c 
	JOIN category child ON child.parent_slug = c.slug
	JOIN filter f ON f.category_slug = c.slug
	WHERE c.slug = $1
	GROUP BY c.slug
`

const categoryCreateQ = `
	INSERT INTO category (slug, title, src, alt, parent_slug) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING slug
`

const categoryUpdateQ = `
	UPDATE category 
	SET title = $1, src = $2, alt = $3, parent_slug = $4 
	WHERE slug = $5
`

const categoryDeleteQ = `
	DELETE FROM category WHERE slug = $1
`

const filterCreateQ = `
	INSERT INTO filter (name, values, filter_type, min_value, max_value, category_slug) 
	VALUES ($1, $2, $3, $4, $5, $6)
`

const filterDeleteQ = `
	DELETE FROM filter WHERE id = ? AND category_slug = ?
`

const itemCategoryDeleteByCategoryQ = `
	DELETE FROM item_category WHERE category_slug = $1
`
