package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	TASKS = "tasks"
)

type Task struct {
	Id        bson.ObjectId `bson:"_id"`
	Completed bool          `bson:"is_completed"`
	Text      string        `bson:"text"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

func GetTaskList(db *mgo.Database) *TaskList {
	tasks := make([]Task, 0)

	collection := db.C(TASKS)
	collection.Find(nil).All(&tasks)

	return &TaskList{tasks}
}

func GetTask(db *mgo.Database, id string) *Task {
	return &Task{}
}
