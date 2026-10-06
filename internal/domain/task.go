package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `json: "id" db: "id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Status      bool       `json: "status" db: "status"`
	CreateAt    time.Time  `json: "create_at" db: "create_at"`
	CompletedAt *time.Time `json: "completed_at" db: "completed_at"`
	Version     int        `json: "version" db: "version`
}

func NewTask(title string, description string) (Task, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:          id,
		Title:       title,
		Description: description,

		Status:      false,
		CreateAt:    time.Now(),
		CompletedAt: nil,
		Version:     0,
	}, nil
}

func (t *Task) SetCompletedStatus(status bool) error {

	if t.Status == status {
		return ErrTaskConflict
	}

	if status {
		completeTime := time.Now()
		t.CompletedAt = &completeTime
	} else {
		t.CompletedAt = nil
	}

	t.Status = status
	t.Version += 1

	return nil
}
