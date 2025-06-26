-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: GetChirp :one
Select * FROM chirps
WHERE $1 = id;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsFromUserID :many
select c.* 
from chirps c
inner join users u ON c.user_id = u.id
where u.id = $1
order by c.created_at asc;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;
