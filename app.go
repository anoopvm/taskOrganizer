package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func NewApp() *App {
	var a App
	return &a
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
	var err error

	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/task/{id:[0-9]+}", a.GetTask).Methods("GET")
	a.Router.HandleFunc("/task", a.CreateTask).Methods("POST")
	a.Router.HandleFunc("/tasks", a.ListTasks).Methods("GET")

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, StatusCode int, message string) {
	respondWithJson(w, StatusCode, message)
}

func respondWithJson(w http.ResponseWriter, StatusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	w.Write(response)
}

func (a *App) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	t := NewTask()
	//Doubtful: interface object t
	t.ID = id
	if err := t.Get(a.DB); err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		//TODO: check for db error and return internal server error
		return
	}
	respondWithJson(w, http.StatusOK, &t)
}

func (a *App) CreateTask(w http.ResponseWriter, r *http.Request) {
	t := NewTask()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	if err = t.Create(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, &t)
}

func (a *App) ListTasks(w http.ResponseWriter, r *http.Request) {
	t := NewTask()
	TaskList, err := t.List(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJson(w, http.StatusOK, TaskList)
}
