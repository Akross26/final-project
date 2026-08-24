package api

import (
	"net/http"

	"github.com/akross26/final-project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJson(w, task)
}
