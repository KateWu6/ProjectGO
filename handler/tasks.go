package handler

import (
	"Project/internal/database"
	"Project/internal/grop"
	"Project/internal/models"
	"Project/internal/my_sort"
	"Project/internal/services"
	"html/template"
	"log"
	"net/http"
	"time"
)

func GetUserTasks(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	usernameInterface := session.Values["username"]
	if usernameInterface == nil {
		log.Println("Пользователь не залогинен, перенаправляем на страницу логина")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, ok := usernameInterface.(string) // Пытаемся привести к строке и проверяем успех
	if !ok {
		log.Println("Ошибка приведения имени пользователя к строке")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Ошибка подключения к БД", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT id_user, energy, lvl, exp FROM users WHERE user_name=$1", username).Scan(&user.ID_user, &user.Energy, &user.Lvl, &user.Exp)
	if err != nil {
		log.Println("Ошибка получения информации о пользователе", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	var tasks []models.Task
	rows, err := db.Query(`SELECT id_task, task_name, task_description, time, date, done, exp, energy_costs 
                          FROM tasks WHERE id_user = $1`, user.ID_user)
	if err != nil {
		log.Println("Ошибка выполнения запроса", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID_task, &task.Task_name, &task.Task_description, &task.Time, &task.Date, &task.Done, &task.Exp, &task.Energy)
		if err != nil {
			log.Println("Ошибка сканирования строки", err)
			continue
		}

		if task.Time.IsZero() {
			log.Println("Внимание: Задача с ID ", task.ID_task, " имеет нулевое время. Проверьте данные в БД.")
			task.Time = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		if task.Date.IsZero() {
			log.Println("Внимание: Задача с ID ", task.ID_task, " имеет нулевую дату. Проверьте данные в БД.")
			task.Date = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Println("Ошибка при итерации по строкам:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Сортируем задачи по сроку завершения
	sortedTasks := my_sort.SortTasksByDeadline(tasks)

	// Группируем задачи по дням
	groupedTasks := grop.GroupTasksByDay(sortedTasks)

	expToNextLevel := services.GetExpForLevel(user.Lvl + 1)

	data := models.TasksPageData{
		Title:          "Задачи",
		IsTasks:        true,
		IsProfile:      false,
		IsAchievements: false,
		Level:          user.Lvl,
		Exp:            user.Exp,
		ExpToNextLevel: int(expToNextLevel),
		Energy:         user.Energy,
		MaxEnergy:      MaxEnergy,
		EnergyUser:     user.Energy,
		GroupedTasks:   groupedTasks,
		ActiveUser:     username,
	}

	tmpl := template.Must(template.New("tasks").
		Funcs(template.FuncMap{
			"groupTasksByDay": grop.GroupTasksByDay,
			"formatDate":      formatDate,
			"formatTime":      formatTime,
		}).
		ParseFiles("templates/tasks.html"))

	err = tmpl.ExecuteTemplate(w, "tasks", data)
	if err != nil {
		log.Println("Ошибка исполнения шаблона", err)
		return
	}
}

// formatDate форматирует дату в строку "YYYY-MM-DD".
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// formatTime форматирует время в строку "HH:MM".
func formatTime(t time.Time) string {
	return t.Format("15:04")
}
