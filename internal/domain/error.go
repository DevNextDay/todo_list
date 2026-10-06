package domain

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrCompletedTask = errors.New("task already completed")
var ErrTaskConflict = errors.New("conflict")
