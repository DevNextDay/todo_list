package dto

import "errors"

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t TaskDTO) ValidateMainTaskInfo() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}

	if t.Description == "" {
		return errors.New("description is empty")
	}

	return nil
}

type CompleteDTO struct {
	Status bool `json:"status"`
}
