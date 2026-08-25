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
const tasksLimit = 50

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	search := r.FormValue("search")

	var (
		tasks []*db.Task
		err   error
	)

	switch {
	case search == "":
		tasks, err = db.Tasks(tasksLimit)

	case isDate(search):
		date, _ := time.Parse(searchDateFormat, search)
		tasks, err = db.TasksByDate(date.Format(dateFormat), tasksLimit)

	default:
		tasks, err = db.TasksSearch(search, tasksLimit)
	}

	if err != nil {
		writeError(w, dbStatus(err), err.Error())
		return
	}

	writeJson(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

func isDate(s string) bool {
	_, err := time.Parse(searchDateFormat, s)
	return err == nil
}
