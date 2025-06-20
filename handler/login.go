package handler

import (
	"Project/internal/my_check"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	// Проверяем существование пользователя
	user, err := my_check.CheckUserExists(username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Пользователь не найден"))
			return
		}
		log.Println("Ошибка обращения к базе данных:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Пользователь не найден"))
		return
	}

	// Проверяем введенный пароль против сохраненного хэша в базе данных
	match, err := ComparePasswordWithHash(password, user.Password)
	if err != nil || !match {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Неправильный пароль."))
		return
	}

	// Проверили пароль, теперь запомним пользователя в сессии
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	session.Values["username"] = username
	err = session.Save(r, w)
	if err != nil {
		log.Println("Ошибка сохранения сессии:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Авторизация успешна, перенаправляем на домашнюю страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ComparePasswordWithHash(password string, hashedPasswd []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPasswd, []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
