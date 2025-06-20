package handler

import (
	"Project/internal/database"
	"Project/internal/models"
	"Project/internal/services"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
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
	err = db.QueryRow("SELECT id_user, energy, lvl, exp FROM users WHERE user_name =$1", username).
		Scan(&user.ID_user, &user.Energy, &user.Lvl, &user.Exp)
	if err != nil {
		log.Println("Ошибка получения информации о пользователе", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Создаем структуру PageData и передаем данные пользователя
	expToNextLevel := services.GetExpForLevel(user.Lvl + 1)

	data := models.CreateTaskPageData{
		Title:          "Создать задачу",
		Level:          user.Lvl,
		Exp:            user.Exp,
		ExpToNextLevel: int(expToNextLevel),
		Energy:         user.Energy,
		MaxEnergy:      MaxEnergy,
		EnergyUser:     user.Energy,
		ActiveUser:     username,
	}

	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/create_task.html", "templates/header.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "create_task", data)
		if err != nil {
			log.Println("Ошибка", err)
			return
		}
		return
	}

	// Если метод POST, обрабатываем форму
	r.ParseForm()
	taskName := strings.TrimSpace(r.FormValue("task_name"))
	taskDescription := strings.TrimSpace(r.FormValue("task_description"))
	date := strings.TrimSpace(r.FormValue("date"))
	time := strings.TrimSpace(r.FormValue("time"))
	expStr := strings.TrimSpace(r.FormValue("exp"))
	energyStr := strings.TrimSpace(r.FormValue("energy_costs"))

	// Преобразование полей опыта и затрат энергии в числа
	exp, err := strconv.Atoi(expStr)
	if err != nil || exp < 0 {
		log.Println("Ошибка преобразования поля EXP в число", err)
		return
	}

	energy, err := strconv.Atoi(energyStr)
	if err != nil || energy < 0 {
		log.Println("Ошибка преобразования поля Energy Costs в число", err)
		http.Error(w, "Некорректный формат затрат энергии", http.StatusBadRequest)
		return
	}

	// Заполняем остальные поля
	done := false

	// SQL-запрос для добавления задачи
	insertQuery := ` INSERT INTO tasks(id_user, task_name, task_description, time, date, done, exp, energy_costs) VALUES ($1, $2, $3, $4::TIME, $5, $6, $7, $8) RETURNING id_task `

	var newTaskID int16
	err = db.QueryRow(insertQuery,
		user.ID_user,
		taskName,
		taskDescription,
		time,
		date,
		done,
		exp,
		energy,
	).Scan(&newTaskID)

	if err != nil {
		log.Println("Ошибка вставки задачи в БД", err)
		http.Error(w, "Ошибка при создании задачи", http.StatusInternalServerError)
		return
	}

	// Уведомление об успешном создании задачи
	fmt.Printf("Новая задача успешно создана с ID=%d\n", newTaskID)

	// Перенаправляем обратно на список задач
	http.Redirect(w, r, "/tasks/", http.StatusSeeOther)
}
