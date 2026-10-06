package http

import (
	"errors"
	"fmt"
	"http_server/internal/transport/http/handler"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *handler.HTTPHandlers
}

func NewHTTPServer(httpHandlers *handler.HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) StartServer() error {

	httpServerPort := fmt.Sprintf(":%s", os.Getenv("HTTP_SERVER_PORT"))

	router := mux.NewRouter()

	router.Path("/tasks").Methods("POST").HandlerFunc(s.httpHandlers.HandlerCreateTask)
	router.Path("/tasks/{id}").Methods("GET").HandlerFunc(s.httpHandlers.HandlerGetTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.httpHandlers.HandlerGetAllTasks)
	router.Path("/tasks/{id}").Methods("PATCH").HandlerFunc(s.httpHandlers.HandlerSetCompletedStatusTask)
	router.Path("/tasks/{id}").Methods("DELETE").HandlerFunc(s.httpHandlers.HandlerDeleteTask)

	if err := http.ListenAndServe(httpServerPort, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
