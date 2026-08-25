package api

import (
	"net/http"
	"time"

	"github.com/akross26/final-project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "метод не поддерживается")
		return
	}

	id := r.FormValue("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, dbStatus(err), err.Error())
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, dbStatus(err), err.Error())
			return
		}
		writeJson(w, http.StatusOK, struct{}{})
		return
	}

	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeError(w, dbStatus(err), err.Error())
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
