package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if date.After(now) {
		return true
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("не могу распарсить дату dstart")
	}

	parts := strings.Split(repeat, " ")

	if parts[0] == "d" {

		if len(parts) < 2 {
			return "", errors.New("не указано число дней")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("число дней не число")
		}

		if days > 400 || days <= 0 {
			return "", errors.New("число дней больше 400 или меньше 1")
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format(dateFormat), nil
	}

	if parts[0] == "y" {

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}

		return date.Format(dateFormat), nil
	}

	if parts[0] == "w" {

		if len(parts) < 2 {
			return "", errors.New("не указаны дни недели")
		}

		daysStr := strings.Split(parts[1], ",")
		var weekdays []int

		for _, s := range daysStr {
			n, err := strconv.Atoi(s)
			if err != nil {
				return "", errors.New("день недели не число")
			}
			if n < 1 || n > 7 {
				return "", errors.New("день недели должен быть от 1 до 7")
			}
			weekdays = append(weekdays, n)
		}

		for {
			date = date.AddDate(0, 0, 1)

			wd := int(date.Weekday())
			if wd == 0 {
				wd = 7
			}

			found := false
			for _, w := range weekdays {
				if w == wd {
					found = true
				}
			}

			if found && afterNow(date, now) {
				break
			}
		}

		return date.Format(dateFormat), nil
	}

	if parts[0] == "m" {

		if len(parts) < 2 {
			return "", errors.New("не указаны дни месяца")
		}

		daysStr := strings.Split(parts[1], ",")
		var mdays []int

		for _, s := range daysStr {
			n, err := strconv.Atoi(s)
			if err != nil {
				return "", errors.New("день месяца не число")
			}
			if n == 0 || n < -2 || n > 31 {
				return "", errors.New("недопустимый день месяца")
			}
			mdays = append(mdays, n)
		}

		var months []int
		if len(parts) >= 3 {
			monthsStr := strings.Split(parts[2], ",")
			for _, s := range monthsStr {
				n, err := strconv.Atoi(s)
				if err != nil {
					return "", errors.New("месяц не число")
				}
				if n < 1 || n > 12 {
					return "", errors.New("недопустимый месяц")
				}
				months = append(months, n)
			}
		}

		for {
			date = date.AddDate(0, 0, 1)

			if len(months) > 0 {
				monthOk := false
				for _, m := range months {
					if int(date.Month()) == m {
						monthOk = true
					}
				}
				if !monthOk {
					continue
				}
			}

			firstNextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
			lastDay := firstNextMonth.AddDate(0, 0, -1).Day()

			dayOk := false
			for _, d := range mdays {
				if d == -1 && date.Day() == lastDay {
					dayOk = true
				} else if d == -2 && date.Day() == lastDay-1 {
					dayOk = true
				} else if d == date.Day() {
					dayOk = true
				}
			}

			if dayOk && afterNow(date, now) {
				break
			}
		}

		return date.Format(dateFormat), nil
	}

	return "", errors.New("неизвестный формат правила")
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "неверный формат now", 400)
			return
		}
	}

	result, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte(result))
}
