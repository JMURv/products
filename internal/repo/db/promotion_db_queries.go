package db

const (
	promoCountQ       = `SELECT COUNT(*) FROM promotion;`
	promoSearchCountQ = `SELECT COUNT(*) FROM promotion WHERE title ILIKE $1 OR slug ILIKE $2;`
)

const promoSearchQ = `
	SELECT slug, title, description, src, alt, lasts_to
	FROM promotion
	WHERE title ILIKE $1 OR slug ILIKE $2
	OFFSET $3 LIMIT $4;
`

const promoCountItemsQ = `
	SELECT COUNT(*) FROM promotion_items WHERE promotion_slug=$1;
`

const promoItemsQ = `
	SELECT 
	    pi.discount, pi.promotion_slug, pi.item_id, 
	    i.title, i.price, i.src, i.alt
	FROM promotion_items pi
	JOIN item i ON i.id = pi.item_id
	WHERE promotion_slug=$1
	GROUP BY pi.promotion_slug
	OFFSET $2 LIMIT $3;
`

const promoListQ = `
	SELECT slug, title, description, src, alt, lasts_to 
	FROM promotion 
	OFFSET $1 LIMIT $2;
`

const promoGetQ = `
	SELECT 
	    slug, title, description, src, alt, lasts_to 
	FROM promotion
	WHERE slug=$1;
`

const promoCreateQ = `
	INSERT INTO promotion (slug, title, description, src, alt, lasts_to) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING slug;
`

const promoItemListQ = `
	SELECT 
	    discount, promotion_slug, item_id 
	FROM promotion_items 
	WHERE promotion_slug = $1
`

const promoItemCreateQ = `
	INSERT INTO promotion_items (discount, promotion_slug, item_id) 
	VALUES ($1, $2, $3);
`

const promoItemUpdateQ = `
	UPDATE promotion_items 
	SET discount = ?, 
	    promotion_slug = ?, 
	    item_id = ?
	WHERE promotion_slug = ?;
`

const promoUpdateQ = `
	UPDATE promotion 
	SET title = $1, description = $2, src = $3, alt = $4, lasts_to = $5 
	WHERE slug = $6;
`

const promoDeleteQ = `
	DELETE FROM promotion WHERE slug = $1
`

const promoItemDelete = `
	DELETE FROM promotion_items WHERE item_id = ? AND promotion_slug = ?;
`
