-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    hashed_password = CASE
        WHEN $1 = TRUE THEN $2
        ELSE hashed_password
    END,
    full_name = CASE 
        WHEN $3 = TRUE THEN $4
        ELSE full_name
    END,
    email = CASE
        WHEN $5 = TRUE THEN $6
        ELSE email
    END
WHERE username = $4
RETURNING *;
