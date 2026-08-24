package api

import (
	"net/http"
	"time"

	"github.com/akross26/final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const searchDateFormat = "02.01.2006"

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	switch {
	case search == "":
		tasks, err = db.Tasks(50)

	case isDate(search):
		date, _ := time.Parse(searchDateFormat, search)
		tasks, err = db.TasksByDate(date.Format(dateFormat), 50)

	default:
		tasks, err = db.TasksSearch(search, 50)
	}

	if err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func isDate(s string) bool {
	_, err := time.Parse(searchDateFormat, s)
	return err == nil
}
