package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"http_server/internal/domain"
	"http_server/internal/service"
	"http_server/internal/transport/http/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	svc *service.TaskService
}

func NewHTTPHandlers(svc *service.TaskService) *HTTPHandlers {
	return &HTTPHandlers{
		svc: svc,
	}
}

/*
pattern: /tasks
method: POST
info: JSON in request body

succeed:
  - status code: 201 Created
  - response body: JSON represent created task

failed:
  - status code: 400, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerCreateTask(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var taskDTO dto.TaskDTO

	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateMainTaskInfo(); err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	task, err := domain.NewTask(taskDTO.Title, taskDTO.Description)
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if err := h.svc.Repo.Create(ctx, task); err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(task, "", "    ")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)

	if _, err := w.Write(b); err != nil {
		fmt.Print("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: GET
info: from pattern

succeed:
  - status code: 200 OK
  - response body: JSON represent found task

failed:
  - status code: 400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerGetTask(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	task, err := h.svc.Repo.GetById(ctx, id)
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		fmt.Print("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks
method: GET
info: -

succeed:
  - status code: 200 OK
  - response body: JSON represent found tasks

failed:
  - status code: 400, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerGetAllTasks(w http.ResponseWriter, r *http.Request) {

	var err error
	var tasks []domain.Task

	q := r.URL.Query()

	if q.Get("status") != "" {
		tasks, err = h.getAllTasksByCompletedStatus(w, r)
	} else {
		tasks, err = h.getAllTasks(w, r)
	}

	if err != nil {
		return
	}

	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Print("failed to write http response:", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: PATCH
info: from pattern + JSON in request body

succeed:
  - status code: 200 OK
  - response body: JSON represent changed task

failed:
  - status code: 400, 404, 409, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerSetCompletedStatusTask(w http.ResponseWriter, r *http.Request) {

	var err error
	var completeDTO dto.CompleteDTO

	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	task, err := h.svc.Repo.GetById(ctx, id)
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())

		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	if err := task.SetCompletedStatus(completeDTO.Status); err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusConflict)
		return
	}

	if err := h.svc.Repo.Update(ctx, task); err != nil {

		errDTO := dto.NewErrorDTO(err.Error(), time.Now())

		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Print("failed to write http response")
		return
	}
}

/*
pattern: /tasks/{title}
method: DELETE
info: from pattern

succeed:
  - status code: 204 No content
  - response body: -

failed:
  - status code: 400, 404, 500, ...
  - response body: JSON with error + time
*/
func (h *HTTPHandlers) HandlerDeleteTask(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return
	}

	if err := h.svc.Repo.Delete(ctx, id); err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		if errors.Is(err, domain.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandlers) getAllTasks(w http.ResponseWriter, r *http.Request) ([]domain.Task, error) {
	ctx := r.Context()
	tasks, err := h.svc.Repo.List(ctx)
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return make([]domain.Task, 0), err
	}
	return tasks, nil
}

func (h *HTTPHandlers) getAllTasksByCompletedStatus(w http.ResponseWriter, r *http.Request) ([]domain.Task, error) {

	ctx := r.Context()
	queryParams := r.URL.Query()

	status, err := strconv.ParseBool(queryParams.Get("status"))
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return make([]domain.Task, 0), err
	}

	tasks, err := h.svc.Repo.ListByCompletedStatus(ctx, status)
	if err != nil {
		errDTO := dto.NewErrorDTO(err.Error(), time.Now())
		http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		return make([]domain.Task, 0), err
	}

	return tasks, nil
}
