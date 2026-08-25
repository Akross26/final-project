package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/akross26/final-project/pkg/db"
)

type taskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func writeJson(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJson: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJson(w, code, taskResponse{Error: msg})
}

func dbStatus(err error) int {
	if err != nil && err.Error() == "задача не найдена" {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("дата представлена в неверном формате")
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("правило повторения укано в неправильном формате")
		}

		if afterNow(now, t) {
			task.Date = next
		}

	} else {
		if afterNow(now, t) {
			task.Date = now.Format(dateFormat)
		}
	}

	return nil

}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	if err := validateTask(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, dbStatus(err), err.Error())
		return
	}

	writeJson(w, http.StatusOK, taskResponse{ID: strconv.FormatInt(id, 10)})
}

func validateTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("не указан заголовок задачи")
	}
	return checkDate(task)
}
