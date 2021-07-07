package main

import (
	"database/sql"
	"fmt"
)

type TaskStatus int

const (
	UnComplete TaskStatus = iota
	Completed
)

type Task interface {
	Get()
	Create()
}

type task struct {
	ID       int        `json:"id"`
	TaskName string     `json:"task"`
	Status   TaskStatus `json:"status"`
}

func NewTask() *task {
	var t task
	return &t
}

func (t *task) Get(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT task, status FROM tasks WHERE id=%d", t.ID)
	return db.QueryRow(statement).Scan(&t.TaskName, &t.Status)
}

func (t *task) Create(db *sql.DB) error {
	t.Status = UnComplete
	statement := fmt.Sprintf("INSERT INTO tasks(task, status) VALUE(\"%s\",%d)", t.TaskName, t.Status)
	//TODO: Save id back to the t
	_, err := db.Exec(statement)
	return err
}

func (t *task) List(db *sql.DB) ([]task, error) {
	statement := "SELECT id, task, status FROM tasks"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []task{}

	for rows.Next() {
		t := task{}
		if err := rows.Scan(&t.ID, &t.TaskName, &t.Status); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
