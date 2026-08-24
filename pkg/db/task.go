package db

import "errors"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	return selectTasks(`SELECT id, date, title, comment, repeat 
	                     FROM scheduler 
	                     ORDER BY date 
	                     LIMIT ?`, limit)
}

func TasksSearch(search string, limit int) ([]*Task, error) {
	like := "%" + search + "%"

	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE title LIKE ? OR comment LIKE ? 
	          ORDER BY date 
	          LIMIT ?`

	return selectTasks(query, like, like, limit)
}

func TasksByDate(date string, limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE date = ? 
	          LIMIT ?`

	return selectTasks(query, date, limit)
}

func selectTasks(query string, args ...any) ([]*Task, error) {
	tasks := []*Task{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	task := Task{}

	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          WHERE id = ?`

	row := db.QueryRow(query, id)

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, errors.New("задача не найдена")
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
	          SET date = ?, title = ?, comment = ?, repeat = ? 
	          WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("задача не найдена")
	}

	return nil
}
