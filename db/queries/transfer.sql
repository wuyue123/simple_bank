-- name: CreateTransfer :one
INSERT INTO Transfers (  
    from_account_id,
    to_account_id,
    amount)
VALUES ($1, $2, $3) 
RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT *
FROM transfers
order by id
LIMIT $1 OFFSET $2;