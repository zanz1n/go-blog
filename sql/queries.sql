-- name: GetUserById :one
SELECT * FROM "users" WHERE "id" = $1;

-- name: GetUserByEmail :one
SELECT * FROM "users" WHERE "email" = $1;

-- name: CreateUser :one
INSERT INTO "users" ("id", "username", "email", "password") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetPostById :one
SELECT * FROM "posts" WHERE "id" = $1;

-- name: GetPostByIdWithUser :one
SELECT
    "posts"."id",
    "posts"."created_at",
    "posts"."updated_at",
    "posts"."title",
    "posts"."content",
    "posts"."topics",
    "posts"."description",
    "posts"."thumb_image",
    "posts"."user_id",
    "users"."created_at" AS "user_created_at",
    "users"."updated_at" AS "user_updated_at",
    "users"."email" AS "user_email",
    "users"."username" AS "user_username",
    "users"."role" AS "user_role"
FROM "posts"
    INNER JOIN "users" ON "posts"."user_id" = "users"."id"
WHERE "posts"."id" = $1;

-- name: GetPartialPostById :one
SELECT
    "posts"."id",
    "posts"."created_at",
    "posts"."updated_at",
    "posts"."title",
    "posts"."description",
    "posts"."thumb_image",
    "posts"."user_id",
    "users"."username" AS "user_username"
FROM "posts"
    INNER JOIN "users" ON "posts"."user_id" = "users"."id"
WHERE "posts"."id" = $1;

-- name: CreatePost :one
INSERT INTO "posts" ("id", "title", "content", "topics", "description", "thumb_image", "user_id") VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;
