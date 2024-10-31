package db

const itemCountQ = `SELECT COUNT(*) FROM item`
const itemCountHitQ = `SELECT COUNT(*) FROM item WHERE is_hit = TRUE;`
const itemCountRecQ = `SELECT COUNT(*) FROM item WHERE is_rec = TRUE;`
const itemCountAttrsQ = `SELECT COUNT(*) FROM item_attr WHERE name ILIKE ?;`

const itemSearchAttrQ = `SELECT id, name, value FROM item_attr WHERE name ILIKE ? OFFSET ? LIMIT ?`

const itemListQ = `SELECT 
    	 i.id, 
    	 i.title,
    	 i.article,
    	 i.description,
    	 i.price,
    	 i.src,
    	 i.alt,
    	 i.in_stock,
    	 i.is_hit,
    	 i.is_rec,
    	 i.parent_id,
    	 i.created_at,
    	 i.updated_at,
    	 ARRAY_AGG(c.title || '|' || c.slug) AS categories
		 FROM item i
		 JOIN item_category ic ON i.id = ic.item_id
		 JOIN category c ON c.slug = ic.category_slug
		 GROUP BY i.id, i.created_at
		 ORDER BY i.created_at DESC 
		 OFFSET ? LIMIT ?`

const itemAttrList = `
	SELECT ia.id, ia.name, ia.value
	FROM item_attr ia
	WHERE ia.item_id = $1
`

const itemMediaList = `
	SELECT im.id, im.src, im.alt
	FROM item_media im
	WHERE im.item_id = $1
`

const itemCategoriesListQ = `
	SELECT ic.item_id, ic.category_slug
	FROM item_category ic
	WHERE ic.item_id = $1
`

const itemRelatedProductList = `
	SELECT rp.item_id, rp.related_item_id
	FROM related_product rp
	WHERE rp.item_id = $1
`

const itemListVarsQ = `SELECT 
	i.id,
	i.title,
	i.article,
	i.description,
	i.price,
	i.src,
	i.alt,
	i.in_stock,
	i.is_hit,
	i.is_rec,
	i.parent_id,
	ARRAY_AGG(im.src || '|' || im.alt) AS media,
	ARRAY_AGG(ia.name || '|' || ia.value) AS attrs
	FROM item i
	JOIN item_attr ia ON ia.item_id = i.id
	JOIN item_media im ON im.item_id = i.id
	WHERE parent_id = $1
	GROUP BY i.id;
`

const itemListRelated = `SELECT 
    i.id, i.title, i.price, i.src, i.alt 
	FROM related_product rp
	JOIN item i ON i.id = rp.item_id
	WHERE item_id = $1;
`

const itemListHitQ = `
	SELECT i.id, i.title, i.price, i.src, i.alt, i.is_hit, i.is_rec
	FROM item i
	WHERE is_hit = TRUE
	OFFSET $1 LIMIT $2;
`

const itemListRecQ = `
	SELECT i.id, i.title, i.price, i.src, i.alt, i.is_hit, i.is_rec
	FROM item i
	WHERE is_rec = TRUE
	OFFSET $1 LIMIT $2;
`

const itemGetByIDQ = `SELECT 
		 i.id, 
		 i.title, 
		 i.article, 
		 i.description,
		 i.price,
		 i.src,
		 i.alt,
		 i.in_stock,
		 i.is_hit,
		 i.is_rec,
		 i.parent_id,
		 i.created_at,
		 i.updated_at,
		 ARRAY_AGG(im.src || '|' || im.alt) AS media,
		 ARRAY_AGG(ia.name || '|' || ia.value) AS attrs,
		 ARRAY_AGG(c.title || '|' || c.slug) AS categories
		 FROM item i
		 JOIN item_media im ON i.id = im.item_id
		 JOIN item_attr ia ON i.id = ia.item_id
		 JOIN item_category ic ON i.id = ic.item_id
		 JOIN category c ON c.slug = ic.category_slug
		 WHERE i.id=?
		 GROUP BY i.id
`

const itemCreateQ = `INSERT INTO item 
    (
     title, 
     article, 
     description, 
     price, 
     src, 
     alt, 
     in_stock,
     is_hit,
     is_rec,
     parent_id
     )
	 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	 RETURNING id;
`

const itemMediaCreateQ = `INSERT INTO item_media (item_id, src, alt) VALUES ($1, $2, $3)`
const itemAttrCreateQ = `INSERT INTO item_attr (item_id, name, value) VALUES ($1, $2, $3)`
const itemCategoryCreateQ = `INSERT INTO item_category (item_id, category_slug) VALUES ($1, $2)`
const itemRelatedProductCreateQ = `INSERT INTO related_product (item_id, related_item_id) VALUES ($1, $2)`

const itemUpdateQ = `UPDATE item 
	SET title = ?, 
	article = ?, 
	description = ?,
	price = ?, 
	src = ?, 
	alt = ?,
	in_stock = ?,
	is_hit = ?,
	is_rec = ?,
	parent_id = ? 
	WHERE id = ?
`

const itemAttrUpdateQ = `
	UPDATE item_attr 
	SET name = ?, value = ? 
	WHERE id = ?
`

const itemDeleteQ = `DELETE FROM item WHERE id = ?`
const itemDeleteByParentQ = `DELETE FROM item WHERE id = ? AND parent_id = ?`
const itemMediaDeleteQ = `DELETE FROM item_media WHERE id = ?`
const itemAttrDeleteQ = `DELETE FROM item_attr WHERE name = ? AND item_id = ?`
const itemCategoryDeleteQ = `DELETE FROM item_category WHERE item_id = ? AND category_slug = ?`
const itemRelatedProductDeleteQ = `DELETE FROM related_product WHERE item_id = ? AND related_item_id = ?`
