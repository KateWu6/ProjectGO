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

func UpdateTask(w http.ResponseWriter, r *http.Request) {

	db, _ := bd.Connect()

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 16)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask bd.Task
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
			exp = $7
			energy_costs = $8
        WHERE id_task = $9
    `

	res, err := db.Exec(query,
		updatedTask.ID_user,
		updatedTask.Task_name,
		updatedTask.Task_description,
		updatedTask.Time,
		updatedTask.Date,
		updatedTask.Done,
		updatedTask.Exp,
		updatedTask.Energy,
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

	w.WriteHeader(http.StatusNoContent)
}
