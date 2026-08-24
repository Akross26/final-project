package api

import (
	"net/http"

	"github.com/akross26/final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJson(w, taskResponse{Error: "не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	writeJson(w, struct{}{})
}
