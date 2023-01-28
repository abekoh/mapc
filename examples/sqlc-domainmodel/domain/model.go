package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserID uuid.UUID

	User struct {
		ID   UserID
		Name string
	}
)

type (
	TaskID uuid.UUID

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
	SubTaskID uuid.UUID

	SubTask struct {
		ID           SubTaskID
		TaskID       TaskID
		UserID       UserID
		Title        string
		Description  string
		RegisteredAt time.Time
	}
)
