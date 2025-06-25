-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: UpdateEmailAndPassword :exec
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3;

-- name: UpgradeUserRed :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE $1 = id
RETURNING *;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE $1 = email;

-- name: GetUserFromID :one
SELECT * FROM users
WHERE $1 = id;

-- name: DeleteAllUsers :exec
DELETE FROM users;
