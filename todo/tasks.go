package todo

import (
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
)

const (
	COLLECTION = "tasks"
)

type Task struct {
	Id        bson.ObjectId `bson:"_id"`
	Completed bool          `bson:"is_completed"`
	Text      string        `bson:"text"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

type TasksHandler struct {
	PostNotAllowed
	PutNotAllowed
	DeleteNotAllowed
}

func (t TasksHandler) Get(values url.Values) Response {
	db, conn := GetNewConnection()
	defer conn.Close()

	var tasks []Task = make([]Task, 0)

	collection := db.C(COLLECTION)

	if err := collection.Find(nil).All(&tasks); err != nil {
		return ResponseError()
	}

	list := TaskList{tasks}
	return Response{http.StatusOK, list}
}
