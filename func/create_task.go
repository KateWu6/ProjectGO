package func_go

import (
	"Project/bd"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {

	db, _ := bd.Connect()

	var newTask bd.Task

	// Чтение данных из тела запроса
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Вставка нового задания в базу данных
	query := `INSERT INTO tasks (id_user, task_name, task_description, time, date, done, exp, energy)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id_task`

	err = db.QueryRow(query,
		newTask.ID_user,
		newTask.Task_name,
		newTask.Task_description,
		newTask.Time,
		newTask.Date,
		newTask.Done,
		newTask.Exp,
		newTask.Energy,
	).Scan(&newTask.ID_task)

	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную задачу с новым ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}
