package func_go

import (
	"database/sql"
	"encoding/json"

	"Project/bd"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	db, _ := bd.Connect()

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
	var user bd.User
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
