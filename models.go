package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	TASKS = "tasks"
)

type Task struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Completed   bool          `bson:"ok" json:"ok"`
	Description string        `bson:"desc" json:"desc"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// Returns a ObjectId from a string if it's possible or nil otherwise
func ParseObjectId(id string) bson.ObjectId {
	var oid bson.ObjectId

	if bson.IsObjectIdHex(id) {
		oid = bson.ObjectIdHex(id)
	}

	return oid
}

// Fetchs all Tasks
func GetTaskList(db *mgo.Database) *TaskList {
	tasks := make([]Task, 0)

	collection := db.C(TASKS)
	collection.Find(nil).All(&tasks)

	return &TaskList{tasks}
}

// Fetch a single Task by the id
func GetTaskById(db *mgo.Database, id string) (*Task, error) {
	var task Task

	collection := db.C(TASKS)
	oid := ParseObjectId(id)
	err := collection.FindId(oid).One(&task)

	return &task, err
}

// Insert the Task at the db
func (t *Task) Insert(db *mgo.Database) error {
	collection := db.C(TASKS)
	// Generate the ObjectId
	t.Id := bson.NewObjectId()
	return collection.Insert(t)
}

// Update the Task at the db
func (t *Task) Update(db *mgo.Database) error {
	collection := db.C(TASKS)
	return collection.UpdateId(t.Id, t)
}

// Remove the Task from db
func (t *Task) Remove(db *mgo.Database) error {
	collection := db.C(TASKS)
	return collection.RemoveId(t.Id)
}
