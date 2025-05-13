-- name: CreateOrderTable :exec
CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id));

-- name: GetAllOrders :many
SELECT * FROM orders;

-- name: GetOrderById :one
SELECT * FROM orders WHERE id = ?;

-- name: CreateOrder :exec
INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?);

-- name: UpdateOrder :exec
UPDATE orders SET price = ?, tax = ?, final_price = ? WHERE id = ?;

-- name: DeleteOrder :exec
DELETE FROM orders WHERE id = ?;


