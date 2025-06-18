package handler

import (
	"Project/bd"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := bd.Connect()

	if r.Method != "POST" {
		http.ServeFile(w, r, "html_files/login.html")
		return
	}

	// Извлекаем данные из POST-запроса
	name := r.FormValue("name")
	password := r.FormValue("password")

	// Валидация полей
	if len(name) == 0 || len(password) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Проверяем, существует ли пользователь с такими данными
	var storedHash []byte
	err = db.QueryRow(`SELECT password_hash FROM users WHERE user_name=$1`, name).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Ошибка при поиске пользователя", http.StatusInternalServerError)
		return
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword(storedHash, []byte(password))
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Авторизация прошла успешно, генерируем новый session_id
	sessionID := uuid.NewString()

	// Ставим сессию в браузер
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
