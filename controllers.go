package main

import (
	"net/http"
	"net/url"
	"strconv"
)

type TaskListController struct{}

func (t *TaskListController) Get(params url.Values) (int, interface{}) {
	db, conn := GetDatabase()
	defer conn.Close()

	return http.StatusOK, GetTaskList(db)
}

type TaskController struct{}

func (t *TaskController) Get(params url.Values) (int, interface{}) {
	db, conn := GetDatabase()
	defer conn.Close()

	task, err := GetTaskById(db, params.Get("id"))
	if err != nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, task
}

func (t *TaskController) Post(params url.Values) (int, interface{}) {
	db, conn := GetDatabase()
	defer conn.Close()

	task := new(Task)
	task.Completed = false
	task.Description = params.Get("desc")

	if err := task.Insert(db); err != nil {
		return http.StatusBadRequest, nil
	}

	return http.StatusOK, task
}

func (t *TaskController) Put(params url.Values) (int, interface{}) {
	db, conn := GetDatabase()
	defer conn.Close()

	task, err := GetTaskById(db, params.Get("id"))
	if err != nil {
		return http.StatusNotFound, nil
	}

	if ok := params.Get("ok"); ok != "" {
		task.Completed, _ = strconv.ParseBool(ok)
	}
	if desc := params.Get("desc"); desc != "" {
		task.Description = desc
	}
	// Do update
	if err := task.Update(db); err != nil {
		return http.StatusBadRequest, nil
	}

	return http.StatusOK, task
}

func (t *TaskController) Delete(params url.Values) (int, interface{}) {
	db, conn := GetDatabase()
	defer conn.Close()

	task, err := GetTaskById(db, params.Get("id"))
	if err != nil {
		return http.StatusNotFound, nil
	}
	// Do remove
	if err := task.Remove(db); err != nil {
		log.Println(err)
		return http.StatusBadRequest, nil
	}

	return http.StatusAccepted, nil
}
