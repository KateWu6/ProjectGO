package my_sort

import (
	"Project/internal/models"
	"time"
)

// Тип для сортировки задач
type Tasks []models.Task

func (t Tasks) Len() int      { return len(t) }           //возвращает длину
func (t Tasks) Swap(i, j int) { t[i], t[j] = t[j], t[i] } //меняет местами элементы по индексу

func (t Tasks) Less(i, j int) bool { //сравнение
	// Проверка наличия действительных значений даты и времени
	if !t[i].Date.IsZero() && !t[i].Time.IsZero() &&
		!t[j].Date.IsZero() && !t[j].Time.IsZero() {
		return t[i].Date.Add(t[i].Time.Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))).Before(
			t[j].Date.Add(t[j].Time.Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))))
	}

	// По умолчанию сортируем задачи с неопределённой датой/временем последними
	return false
}
