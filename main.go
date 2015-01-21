package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello, Web!")
}

func main() {
	db := GetMongoConnection()
	db.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	defer db.Close()

	log.Println("Starting server...")

	list := new(TaskListController)
	http.HandleFunc("/tasks", RestController(list))

	task := new(TaskListController)
	http.HandleFunc("/task", RestController(task))

	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(":8000", nil)
}
