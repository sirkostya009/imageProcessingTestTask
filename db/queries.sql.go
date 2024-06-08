// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package db

import (
	"context"
)

const createImage = `-- name: CreateImage :exec
insert into images (user_id, url)
values ($1, $2)
`

type CreateImageParams struct {
	UserID int32
	Url    string
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) error {
	_, err := q.db.Exec(ctx, createImage, arg.UserID, arg.Url)
	return err
}

const createUser = `-- name: CreateUser :one
insert into users (username, email, password_hash)
values ($1, $2, $3)
returning id, username, email, password_hash
`

type CreateUserParams struct {
	Username     string
	Email        *string
	PasswordHash string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
	)
	return i, err
}

const deleteOldImages = `-- name: DeleteOldImages :exec
delete from images where user_id = $1
`

func (q *Queries) DeleteOldImages(ctx context.Context, userID int32) error {
	_, err := q.db.Exec(ctx, deleteOldImages, userID)
	return err
}

const getImages = `-- name: GetImages :many
select id, user_id, url from images where user_id = $1
`

func (q *Queries) GetImages(ctx context.Context, userID int32) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImages, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Image
	for rows.Next() {
		var i Image
		if err := rows.Scan(&i.ID, &i.UserID, &i.Url); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByUsername = `-- name: GetUserByUsername :one
select id, username, email, password_hash from users where username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
	)
	return i, err
}

const usernameTaken = `-- name: UsernameTaken :one
select exists(select 1 from users where username = $1)
`

func (q *Queries) UsernameTaken(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRow(ctx, usernameTaken, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}