package main

import (
	"log"
	"net/http"
	"time"
)

type Task struct {
	Id   int       `json: "id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

func New() *TaskStore
func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int
func (ts *TaskStore) GetTask(id int) (Task, error)
func (ts *TaskStore) DeleteTask(id int) error
func (ts *TaskStore) DeleteAllTasks() error
func (ts *TaskStore) GetAllTasks() []Task
func (ts *TaskStore) GetTaskByTag(tag string) []Task
func (ts *TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []Task

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("POST /task/", server.createTaskHandler)
	mux.HandleFunc("GET /task/", server.getAllTasksHandler)
	mux.HandleFunc("DELETE /task/", server.deleteAllTasksHandler)
	mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	mux.HandleFunc("DELETE /task/{id}", server.deleteTaskHandler)
	mux.HandleFunc("GET /tag/{tag}", server.tagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+"4112"), mux)
}
