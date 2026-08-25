package api

import (
	"net/http"

	"github.com/akross26/final-project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJson(w, http.StatusOK, task)
}
