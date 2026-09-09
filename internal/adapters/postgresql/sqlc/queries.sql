-- name: ListProducts :many
SELECT * from products;

-- name: FindProductById :one
SELECT * from products where id = $1;

-- name: CreateOrder :one
INSERT INTO orders (customer_id) values ($1)
RETURNING *;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity, price_cents)
values ($1, $2, $3, $4);

-- name: UpdateProductQuantity :exec
UPDATE products
SET quantity = $2
WHERE id = $1;
