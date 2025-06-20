package my_sort

import (
	"Project/internal/models"
	"sort"
)

// Функция для сортировки списка задач
func SortTasksByDeadline(tasks []models.Task) Tasks {
	sort.Sort(Tasks(tasks))
	return tasks
}
