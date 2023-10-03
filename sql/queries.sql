-- name: GetUserById :one
SELECT * FROM "users" WHERE "id" = $1;

-- name: CreateUser :one
INSERT INTO "users" ("id", "username", "email", "password") VALUES ($1, $2, $3, $4) RETURNING *;
