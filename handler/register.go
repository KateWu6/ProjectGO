package handler

import (
	"Project/internal/database"
	"Project/internal/my_check"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/register.html")
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

	if len(username) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Все поля обязательны для заполнения"))
		return
	}

	isUnique, err := my_check.CheckUsername(username)
	if err != nil {
		log.Println("Ошибка проверки уникальности пользователя:")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if !isUnique {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Такое имя уже существует"))
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //хэширование
	if err != nil {
		log.Println("Ошибка хэширования пароля:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Ошибка подключения к БД", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	achievement := 0
	energy := 10
	lvl := 1
	exp := 0

	insertQuery := `INSERT INTO users (user_name, password, achievement, energy, lvl, exp) VALUES($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(insertQuery, username, hashedPwd, achievement, energy, lvl, exp)
	if err != nil {
		log.Println("Ошибка вставки пользователя в БД", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
