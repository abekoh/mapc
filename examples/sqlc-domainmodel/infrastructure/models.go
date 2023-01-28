// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package infrastructure

import (
	"database/sql"
	"time"
)

type SubTask struct {
	ID           string
	TaskID       string
	UserID       string
	Title        string
	Description  string
	RegisteredAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Task struct {
	ID           string
	UserID       string
	Title        string
	Description  string
	StoryPoint   sql.NullInt64
	RegisteredAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type User struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
