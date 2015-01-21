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

	list := new(TaskListController)
	http.HandleFunc("/tasks", RestController(list))

	task := new(TaskController)
	http.HandleFunc("/task", RestController(task))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Starting server at %s...", port)

	http.HandleFunc("/", helloHandler)
	http.ListenAndServe(port, nil)
}
