package api

import (
	"encoding/json"
	"net/http"

	"github.com/akross26/final-project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор задачи")
		return
	}

	if err := validateTask(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, dbStatus(err), err.Error())
		return
	}

	writeJson(w, http.StatusOK, struct{}{})
}
