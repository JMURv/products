package db

const orderGetQ = `
SELECT 
	o.id, 
	o.status, 
	o.total_amount, 
	o.fio, 
	o.tel, 
	o.email,
	o.address, 
	o.delivery, 
	o.payment_method,
	o.user_id,
	o.created_at,
	o.updated_at,
	ARRAY_AGG(oi.id || '|' || oi.item_id || '|' || oi.quantity) AS order_items
FROM "order" o
JOIN "order_item" oi ON o.id = oi.order_id
WHERE o.id=$1 
`

const userOrderCountQ = `SELECT COUNT(*) FROM "order" WHERE user_id=?`
const userOrdersQ = `
SELECT 
	o.id, 
	o.status, 
	o.total_amount, 
	o.fio, 
	o.tel, 
	o.email,
	o.address, 
	o.delivery, 
	o.payment_method,
	o.user_id,
	o.created_at,
	o.updated_at,
	ARRAY_AGG(oi.id || '|' || oi.item_id || '|' || oi.quantity) AS order_items
FROM "order" o
JOIN "order_item" oi ON o.id = oi.order_id
WHERE user_id=$1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3
`

const orderCreateQ = `
INSERT INTO "order" 
    (status, total_amount, fio, tel, email, address, delivery, payment_method, user_id, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING id
`

const orderUpdateQ = `
UPDATE "order"
SET status=$1, total_amount=$2, fio=$3, tel=$4, email=$5, address=$6, delivery=$7, payment_method=$8, updated_at=NOW()
WHERE id=$9
`

const orderItemListQ = `
SELECT
	oi.id,
	oi.item_id,
	oi.quantity
FROM "order_item" oi
WHERE oi.order_id=$1
`

const orderItemCreateQ = `INSERT INTO "order_item" (order_id, item_id, quantity) VALUES ($1, $2, $3)`

const orderItemDeleteQ = `DELETE FROM "order_item" WHERE id=$1`

const orderCancelQ = `UPDATE "order" SET status=$1, updated_at = NOW() WHERE id=$2`

const orderItemPriceQ = `SELECT price FROM item WHERE id = $1`
