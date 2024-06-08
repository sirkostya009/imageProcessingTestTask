-- name: CreateUser :one
insert into users (username, email, password_hash)
values ($1, $2, $3)
returning *;

-- name: CreateImage :exec
insert into images (user_id, url)
values ($1, $2);

-- name: GetUserByUsername :one
select * from users where username = $1;

-- name: GetImages :many
select * from images where user_id = $1;

-- name: UsernameTaken :one
select exists(select 1 from users where username = $1);

-- name: DeleteOldImages :exec
delete from images where user_id = $1;
