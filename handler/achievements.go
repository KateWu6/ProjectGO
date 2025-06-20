package handler

import (
	"Project/internal/database"
	"Project/internal/models"
	"Project/internal/services"
	"html/template"
	"log"
	"net/http"
)

func GetUserAchievement(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
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
		log.Println("Ошибка подключения к БД")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var user models.User
	err = db.QueryRow("SELECT id_user, energy, lvl, exp FROM users WHERE user_name =$1", username).Scan(&user.ID_user, &user.Energy, &user.Lvl, &user.Exp)
	if err != nil {
		log.Println("Ошибка получения информации о пользователе", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	var achievements []models.Achievement
	rows, err := db.Query(`
	SELECT acievement_name, acievement_description, done
	FROM acievement 
	WHERE id_user = $1`, user.ID_user)
	if err != nil {
		log.Println("Ошибка выполнения запроса", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var achievement models.Achievement
		err = rows.Scan(&achievement.Achievement_name, &achievement.Achievement_description, &achievement.Done)
		if err != nil {
			log.Println("Ошибка сканирования строки", err)
			continue
		}
		achievements = append(achievements, achievement)
	}

	if err = rows.Err(); err != nil {
		log.Println("Ошибка при итерации по строкам:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	expToNextLevel := services.GetExpForLevel(user.Lvl + 1)

	data := models.AchievementPageData{
		Title:          "Достижения",
		IsProfile:      false,
		IsTasks:        false,
		IsAchievements: true,
		Level:          user.Lvl,
		Exp:            user.Exp,
		ExpToNextLevel: int(expToNextLevel),
		Energy:         user.Energy,
		MaxEnergy:      MaxEnergy,
		EnergyUser:     user.Energy,
		Achievements:   achievements,
		ActiveUser:     username,
	}

	tmpl, err := template.ParseFiles("templates/achievements.html")
	if err != nil {
		log.Println("Ошибка разбора шаблона", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "achievements", data)
	if err != nil {
		log.Println("Ошибка исполнения шаблона", err)
		return
	}
}
