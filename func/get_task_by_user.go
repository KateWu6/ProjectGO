package func_go

import (
	"Project/bd"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func GetTasksByUser(w http.ResponseWriter, r *http.Request) {

	db, _ := bd.Connect()

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

	tasks := []bd.Task{}

	for rows.Next() {
		var t bd.Task
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
