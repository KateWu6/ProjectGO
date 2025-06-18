package func_go

import (
	"Project/bd"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	db, _ := bd.Connect()

	var newUser bd.User

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

	newUser.Energy = 10
	newUser.Lvl = 1
	newUser.Achievement = 0

	// Вставляем нового пользователя в базу, исключая ID (предполагается автоинкремент)
	query := `INSERT INTO users (user_name, password_hash, achievement, energy, lvl) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var insertedID int64

	err = db.QueryRow(query, newUser.Name, string(hashedPassword), newUser.Achievement, newUser.Energy, newUser.Lvl).Scan(&insertedID)
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с созданным пользователем
	newUser.ID_user = uint16(insertedID)

	// Не возвращаем пароль
	newUser.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}
