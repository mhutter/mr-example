package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mhutter/mr"
)

type Todo struct {
	mr.Base `json:",inline" bson:",inline"`

	Description string `json:"description" bson:"description"`
}

func main() {
	// Connect to MongoDB
	repo := mr.MustAutoconnect("mr-example")

	// Create router
	r := mux.NewRouter()

	// Use the MongoRepo Middleware
	r.Use(mr.Middleware(repo))

	r.HandleFunc("/api/todos", AddTodo).Methods("POST")
	r.HandleFunc("/api/todos", ListTodos).Methods("GET")

	log.Fatalln(http.ListenAndServe(":3001", r))
}

func AddTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&todo)

	// do some validation here

	// Insert into the DB
	mr.C(r, "todos").Insert(&todo)

	// Return the inserted item (including its ID)
	json.NewEncoder(w).Encode(&todo)
}

func ListTodos(w http.ResponseWriter, r *http.Request) {
	todos := []Todo{}

	// Fetch all ToDos from the DB
	mr.C(r, "todos").FindAll(&todos)

	// return JSON
	json.NewEncoder(w).Encode(&todos)
}
