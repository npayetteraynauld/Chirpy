-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  token, 
  created_at, 
  updated_at,
  user_id,
  expires_at
)
VALUES (
  $1,
  NOW(),
  NOW(),
  $2,
  NOW() + INTERVAL '60 days'
)
RETURNING *;

-- name: GetRefreshTokenFromId :one
SELECT * FROM refresh_tokens
WHERE $1 = token;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users u
INNER JOIN refresh_tokens rf
ON u.id = rf.user_id
WHERE rf.token = $1;

-- name: GetRefreshTokenFromUserID :one
SELECT * FROM refresh_tokens 
WHERE user_id = $1
AND expires_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

