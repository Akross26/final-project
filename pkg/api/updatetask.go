package api

import (
	"encoding/json"
	"net/http"

	"github.com/akross26/final-project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, taskResponse{Error: "ошибка десериализации JSON"})
		return
	}

	if task.ID == "" {
		writeJson(w, taskResponse{Error: "не указан идентификатор задачи"})
		return
	}

	if err := validateTask(&task); err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, taskResponse{Error: err.Error()})
		return
	}

	writeJson(w, struct{}{})
}
