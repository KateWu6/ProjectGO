package handler

import (
	"Project/bd"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	db, err := bd.Connect()

	if r.Method != "POST" {
		http.ServeFile(w, r, "html_files/register.html")
		return
	}

	// Читаем поля из POST-запроса
	name := r.FormValue("name")
	password := r.FormValue("password")

	// Валидация входных данных
	if len(name) == 0 || len(password) == 0 {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка обработки пароля", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базу данных
	query := `INSERT INTO users (user_name, password_hash) VALUES ($1, $2) RETURNING id`
	var insertedID int64

	err = db.QueryRow(query, name, hashedPassword).Scan(&insertedID)
	if err != nil {
		http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
		return
	}

	// Генерируем уникальный session_id
	sessionID := uuid.NewString()

	// Устанавливаем куки с session_id
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})

	// Переход на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
