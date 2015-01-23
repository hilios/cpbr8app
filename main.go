package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	WEBAPP_URL = "http://hilios.github.io/cpbr8app/"
)

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, WEBAPP_URL, http.StatusTemporaryRedirect)
}

func main() {
	db := GetMongoConnection()
	db.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	defer db.Close()

	list := new(TaskListController)
	http.HandleFunc("/tasks", RestController(list))

	task := new(TaskController)
	http.HandleFunc("/task", RestController(task))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Starting server at %s...", port)

	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(port, nil)
}
