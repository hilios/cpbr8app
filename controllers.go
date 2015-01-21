package main

import (
	"net/http"
	"net/url"
)

type TaskListController struct{}

func (t *TaskListController) Get(values url.Values) (int, interface{}) {
	db, conn := GetDatabase()
	defer conn.Close()

	return http.StatusOK, GetTaskList(db)
}

type TaskController struct{}

func (t *TaskController) Get(values url.Values) (int, interface{}) {
	return http.StatusAccepted, nil
}

func (t *TaskController) Post(values url.Values) (int, interface{}) {
	return http.StatusAccepted, nil
}

func (t *TaskController) Put(values url.Values) (int, interface{}) {
	return http.StatusAccepted, nil
}

func (t *TaskController) Delete(values url.Values) (int, interface{}) {
	return http.StatusAccepted, nil
}
