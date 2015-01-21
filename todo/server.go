package todo

import (
	"net/http"
)

var mux *http.ServeMux

func init() {
	mux = http.NewServeMux()
	mux.HandleFunc("/tasks", RestHandler(TasksHandler{}))
}

func GetServerMux() *http.ServeMux {
	return mux
}
