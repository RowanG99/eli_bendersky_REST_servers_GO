package taskstore

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
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

func renderJson(w http.ReponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)

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

	renderJson(w, task)
}

func (ts *taskServer) getAllTasksHandler(w http.ReponseWriter, req *http.Request) {
	log.Printf("handling get all tasks at %s\n", req.URL.Path)

	tasks, err := ts.store.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJson(w, tasks)
}

func (ts *taskServer) getTaskByTagHanlder(w http.ReponseWriter, req *http.Request) {
	log.Printf("handling get task by tag at %s\n", req.URL.Path)

	tag, err := strconv.Atoi(req.PathValue("tag"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	tasks, err := ts.store.GetTaskByTag(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	renderJson(w, tasks)
}

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer()
	mux.HandleFunc("POST /task/", server.createTaskHandler)
	mux.HandleFunc("GET /task/", server.getAllTasksHandler)
	mux.HandleFunc("DELETE /task/", server.deleteAllTasksHandler)
	mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	mux.HandleFunc("DELETE /task/{id}/", server.deleteTaskHandler)
	mux.HandleFunc("GET /tag/{tag}/", server.getTaskByTagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)

	log.Fatal(http.ListenAndServe("localhost:"+"4112"), mux)
}
