package api

import (
	"net/http"
	"time"

	"github.com/akross26/final-project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJson(w, taskResponse{Error: "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, taskResponse{Error: err.Error()})
			return
		}
		writeJson(w, struct{}{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	writeJson(w, struct{}{})
}
