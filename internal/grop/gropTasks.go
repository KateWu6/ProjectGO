package grop

import (
	"Project/internal/models"
)

// Группируем задачи по дню
func GroupTasksByDay(tasks []models.Task) []models.DayGroup {
	groups := make(map[string][]models.Task)
	for _, task := range tasks {
		dateKey := task.Date.Format("2006-01-02") // Форматирование ключа группы
		groups[dateKey] = append(groups[dateKey], task)
	}

	result := make([]models.DayGroup, 0, len(groups)) //срез
	for date, groupedTasks := range groups {
		result = append(result, models.DayGroup{Date: date, Tasks: groupedTasks})
	}

	return result
}
