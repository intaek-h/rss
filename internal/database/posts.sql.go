// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
insert into posts (
    id,
    created_at,
    updated_at,
    title,
    description,
    published_at,
    url,
    feed_id
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning id, created_at, updated_at, title, description, published_at, url, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description sql.NullString
	PublishedAt time.Time
	Url         string
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Description,
		arg.PublishedAt,
		arg.Url,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.PublishedAt,
		&i.Url,
		&i.FeedID,
	)
	return i, err
}
