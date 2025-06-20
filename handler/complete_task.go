package handler

import (
	"Project/internal/database"
	"Project/internal/services"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	taskIdStr := vars["id"]
	taskId, err := strconv.ParseInt(taskIdStr, 10, 16)
	if err != nil {
		http.Error(w, "Неверный ID задачи", http.StatusBadRequest)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Ошибка подключения к БД")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 1. Получаем данные о пользователе и задаче
	var userId uint16
	var taskExp int
	var taskEnergy int
	var taskDone bool

	err = db.QueryRow(`
  	SELECT u.id_user, t.exp, t.energy_costs
  	FROM users u
  	JOIN tasks t ON u.id_user = t.id_user
  	WHERE u.user_name = $1 AND t.id_task = $2`,
		username, taskId).Scan(&userId, &taskExp, &taskEnergy)

	if err == sql.ErrNoRows {
		log.Println("Запись не найдена")
		http.Error(w, "Запись не найдена", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Ошибка получения данных о пользователе и задаче:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// **Проверяем, выполнена ли задача:**
	if taskDone {
		log.Println("Задача уже выполнена")
		w.WriteHeader(http.StatusConflict) // Или другой подходящий статус
		fmt.Fprint(w, "Задача уже выполнена")
		return
	}

	// 2. Обновляем задачу (отмечаем как выполненную)
	_, err = db.Exec("UPDATE tasks SET done = TRUE WHERE id_task = $1", taskId)
	if err != nil {
		log.Println("Ошибка обновления задачи", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// 3. Обновляем опыт и энергию пользователя
	err = services.UpdateUserExpAndEnergy(db, userId, taskExp, taskEnergy) // Вызываем функцию для обновления
	if err != nil {
		log.Println("Ошибка обновления опыта и энергии пользователя", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// 4. Отправляем успешный ответ
	w.WriteHeader(http.StatusOK) // Отправляем код 200 OK
	fmt.Fprint(w, "Задача успешно выполнена")
}
