package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const searchDateFormat = "02.01.2006"

func AddTask(task *Task) (int64, error) {
	if database == nil {
		return 0, fmt.Errorf("database is not initialized")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := database.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return id, nil
}

func Tasks(limit int) ([]*Task, error) {
	if database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := database.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func TasksBySearch(search string, limit int) ([]*Task, error) {
	if database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date`
	rows, err := database.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	all, err := scanTasks(rows)
	if err != nil {
		return nil, err
	}

	needle := strings.ToLower(search)
	result := []*Task{}
	for _, t := range all {
		if strings.Contains(strings.ToLower(t.Title), needle) ||
			strings.Contains(strings.ToLower(t.Comment), needle) {
			result = append(result, t)
			if len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}

func TasksByDate(date string, limit int) ([]*Task, error) {
	if database == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE date = ? ORDER BY date LIMIT ?`
	rows, err := database.Query(query, date, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

func scanTasks(rows *sql.Rows) ([]*Task, error) {
	tasks := []*Task{}

	for rows.Next() {
		task := &Task{}
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tasks, nil
}

func ParseSearchDate(search string) (string, bool) {
	t, err := time.Parse(searchDateFormat, search)
	if err != nil {
		return "", false
	}
	return t.Format("20060102"), true
}
