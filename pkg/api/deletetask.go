package api

import (
	"net/http"

	"github.com/akross26/final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, dbStatus(err), err.Error())
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
