package domain

import (
	"context"
	"time"
)

type (
	UserID string

	User struct {
		ID   UserID
		Name string
	}

	UserRepository interface {
		Get(ctx context.Context, id UserID) (*User, error)
		List(ctx context.Context) (*User, error)
	}
)

func (i UserID) String() string {
	return string(i)
}

type (
	TaskID string

	TaskType int

	Task struct {
		ID           TaskID
		UserID       UserID
		Title        string
		Description  string
		StoryPoint   *int
		RegisteredAt time.Time
	}
)

const (
	Story TaskType = iota
	Kaizen
	Bug
)

type (
	SubTaskID string

	SubTask struct {
		ID           SubTaskID
		TaskID       TaskID
		UserID       UserID
		Title        string
		Description  string
		RegisteredAt time.Time
	}
)
