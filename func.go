package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User

	// Декодируем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Вставляем нового пользователя в базу, исключая ID (предполагается автоинкремент)
	query := `INSERT INTO users (user_name, password_hash, achievement, energy, lvl) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var insertedID int64

	err = db.QueryRow(query, newUser.Name, string(hashedPassword), newUser.Achievement, newUser.Energy, newUser.Lvl).Scan(&insertedID)
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с созданным пользователем (можно оставить только ID или вернуть все данные)
	newUser.ID_user = uint16(insertedID)

	// Не возвращаем пароль
	newUser.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр id из URL, например /user?id=123
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id64, err := strconv.ParseUint(idStr, 10, 16)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}
	id := uint16(id64)

	// Запрос к базе данных
	var user User
	query := `SELECT id, user_name, achievement, energy, lvl FROM users WHERE id = $1`
	err = db.QueryRow(query, id).Scan(
		&user.ID_user,
		&user.Name,
		&user.Achievement,
		&user.Energy,
		&user.Lvl,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// Возвращаем данные пользователя (без пароля)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	// Чтение данных из тела запроса
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Вставка нового задания в базу данных
	query := `INSERT INTO tasks (id_user, task_name, task_description, time, date, done)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_task`

	err = db.QueryRow(query,
		newTask.ID_user,
		newTask.Task_name,
		newTask.Task_description,
		newTask.Time,
		newTask.Date,
		newTask.Done,
	).Scan(&newTask.ID_task)

	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную задачу с новым ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func getTasksByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]

	userId, err := strconv.ParseInt(userIdStr, 10, 16)
	if err != nil {
		http.Error(w, "Invalid userId", http.StatusBadRequest)
		return
	}

	// SQL-запрос для получения задач пользователя
	rows, err := db.Query(`
        SELECT id_task, id_user, task_name, task_description, time, date, done 
        FROM tasks 
        WHERE id_user = $1`, userId)
	if err != nil {
		log.Println("DB query error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID_task,
			&t.ID_user,
			&t.Task_name,
			&t.Task_description,
			&t.Time,
			&t.Date,
			&t.Done,
		)
		if err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		log.Println("Rows iteration error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Проверяем, что id из URL совпадает с id в теле (если указан)
	if updatedTask.ID_task != 0 && updatedTask.ID_task != int16(id) {
		http.Error(w, "Task ID mismatch", http.StatusBadRequest)
		return
	}

	// Обновляем запись в базе
	query := `
        UPDATE tasks SET 
            id_user = $1,
            task_name = $2,
            task_description = $3,
            time = $4,
            date = $5,
            done = $6
        WHERE id_task = $7
    `

	res, err := db.Exec(query,
		updatedTask.ID_user,
		updatedTask.Task_name,
		updatedTask.Task_description,
		updatedTask.Time,
		updatedTask.Date,
		updatedTask.Done,
		id,
	)
	if err != nil {
		log.Println("DB update error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content — обновление прошло успешно, тело не отправляем
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Выполняем удаление задачи по id_task
	res, err := db.Exec("DELETE FROM tasks WHERE id_task = $1", id)
	if err != nil {
		log.Println("DB delete error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content — успешное удаление, тело не отправляем
}
