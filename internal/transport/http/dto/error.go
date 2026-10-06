package dto

import (
	"encoding/json"
	"time"
)

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func NewErrorDTO(message string, time time.Time) ErrorDTO {
	return ErrorDTO{
		Message: message,
		Time:    time,
	}
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
