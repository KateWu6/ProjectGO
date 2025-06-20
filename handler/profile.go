package handler

import (
	"Project/internal/database"
	"Project/internal/models"
	"Project/internal/services"
	"html/template"
	"log"
	"net/http"
)

const (
	MaxEnergy = 100
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Перенаправляем на логин
		return
	}

	usernameInterface := session.Values["username"]
	if usernameInterface == nil {
		log.Println("Пользователь не залогинен, перенаправляем на страницу логина")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, ok := usernameInterface.(string) // Пытаемся привести к строке и проверяем успех
	if !ok {
		log.Println("Ошибка приведения имени пользователя к строке")
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

	var user models.User
	err = db.QueryRow("SELECT id_user, user_name, energy, lvl, exp FROM users WHERE user_name = $1", username).Scan(&user.ID_user, &user.Name, &user.Energy, &user.Lvl, &user.Exp)
	if err != nil {
		log.Println("Ошибка получения информации о пользователе", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	expToNextLevel := services.GetExpForLevel(user.Lvl + 1)

	data := models.PageData{
		Title:          "Профиль",
		IsProfile:      true,
		IsTasks:        false,
		IsAchievements: false,
		Level:          user.Lvl,
		Exp:            user.Exp,
		ExpToNextLevel: int(expToNextLevel),
		Energy:         user.Energy,
		MaxEnergy:      MaxEnergy,
		EnergyUser:     user.Energy,
		User:           user,
	}

	tmpl, err := template.ParseFiles("templates/profile.html", "templates/header.html") // Добавляем header.html
	if err != nil {
		log.Println("Ошибка разбора шаблона:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "profile.html", data) // Используем ExecuteTemplate
	if err != nil {
		log.Println("Ошибка исполнения шаблона:", err)
		return
	}

	log.Println("Шаблон профиля успешно исполнен.")
}
