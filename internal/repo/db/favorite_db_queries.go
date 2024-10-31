package db

const favGetQ = `
	SELECT f.user_id, f.item_id, i.id, i.title, i.src, i.alt, i.price 
	FROM favorites f
	JOIN item i ON i.id = f.item_id
	WHERE user_id = ?
`

const favAddQ = `
	INSERT INTO favorites (user_id, item_id) 
	VALUES (?, ?)
	ON CONFLICT (user_id, item_id) DO NOTHING
	RETURNING item_id
`

const favDelQ = `
	DELETE FROM favorites
	WHERE user_id = ? AND item_id = ?
`
