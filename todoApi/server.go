package main
import "strconv"

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"io/ioutil"
)

type Todo struct {
	id	int	`json:"id"`
	title	string	`json."title"`
}

var TodoList = []Todo{
	{id: 1, title: "Learn Go"},
	{id: 2, title: "Build a web server"},
	{id: 3, title: "Write some code"},
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	loadTask()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TodoList)

}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	fmt.Println("Decoded Todo:", newTodo)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	fmt.Println(newTodo)
	TodoList = append(TodoList, newTodo)
	saveTask()
	fmt.Println("Updated TodoList:", TodoList)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	found := false
	for i, todo := range TodoList {
		if todo.id == id {
			TodoList = append(TodoList[:i], TodoList[i+1:]...)
			found = true
			break
		}
	}

	saveTask()

	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

}

func updateTodoHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var updatedTodo Todo
	err := json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	found := false
	for i, todo := range TodoList {
		if todo.id == updatedTodo.id {
			TodoList[i].title = updatedTodo.title
			found = true
			break
		}
	}
	saveTask()
	if !found {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTodo)
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Logic BEFORE the actual handler
		fmt.Printf("Started %s %s\n", r.Method, r.URL.Path)

		// 2. Call the "next" handler in the chain
		next(w, r)

		// 3. Logic AFTER the actual handler
		fmt.Printf("Completed in %v\n", time.Since(start))
	}
}

const fileName = "todo.json"

func saveTask(){
	// Convert slice to JSON

	data, err := json.MarshalIndent(TodoList, "", " ")

	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}
	// Write JSON data to file
	ioutil.WriteFile(fileName, data, 0644)
}

func loadTask() {
	// Read JSON data from file
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	json.Unmarshal(data, &TodoList)
}

func main(){
	http.HandleFunc("/todos", loggerMiddleware(getTodosHandler))
	http.HandleFunc("/add", loggerMiddleware(addTodoHandler))
	http.HandleFunc("/delete", loggerMiddleware(deleteTodo))
	http.HandleFunc("/update", loggerMiddleware(updateTodoHandler))
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}