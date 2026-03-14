package taskstore

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handleing get task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

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
