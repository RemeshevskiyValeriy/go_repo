package cache

import "fmt"

func TaskByIDKey(id int64) string {
	return fmt.Sprintf("tasks:task:%d", id)
}

func TasksListKey(page, limit int) string {
	if page <= 0 || limit <= 0 {
		return "tasks:list"
	}
	return fmt.Sprintf("tasks:list?page=%d&limit=%d", page, limit)
}
