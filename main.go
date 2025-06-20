package main

import (
	"Project/handler"
	"Project/internal/database"
	"Project/internal/services"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handler.ProfilePage).Methods("GET")

	fs := http.FileServer(http.Dir("templates"))
	log.Println("Обслуживаем статические файлы из templates")
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", fs))

	router.HandleFunc("/login", handler.LoginHandler)
	router.HandleFunc("/register", handler.RegisterHandler)
	router.HandleFunc("/logout", handler.LogoutHandler).Methods("POST")

	router.HandleFunc("/tasks/", handler.GetUserTasks)
	router.HandleFunc("/achievements/", handler.GetUserAchievement)

	router.HandleFunc("/tasks/create_task/", handler.CreateTask)

	router.HandleFunc("/tasks/{id}/complete", handler.CompleteTaskHandler).Methods("POST")

	go func() {
		for {
			// Получаем список всех user_id из базы данных.
			db, err := database.Connect()
			if err != nil {
				fmt.Println("Ошибка подключения к базе данных:", err)
				time.Sleep(24 * time.Hour) // Повторная попытка через 24 часа
				continue
			}

			rows, err := db.Query("SELECT id_user FROM users")
			if err != nil {
				fmt.Println("Ошибка получения списка пользователей:", err)
				db.Close()
				time.Sleep(24 * time.Hour) // Повторная попытка через 24 часа
				continue
			}
			defer rows.Close()

			var userIDs []uint16
			for rows.Next() {
				var userID uint16
				if err := rows.Scan(&userID); err != nil {
					fmt.Println("Ошибка сканирования ID пользователя:", err)
					continue
				}
				userIDs = append(userIDs, userID)
			}
			db.Close()

			now := time.Now()
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			nextMidnight := midnight.Add(24 * time.Hour)
			durationUntilMidnight := nextMidnight.Sub(now)

			time.Sleep(durationUntilMidnight)

			for _, userID := range userIDs {
				if err := services.ReplenishEnergy(userID); err != nil { // Вызов через имя пакета
					fmt.Printf("Ошибка пополнения энергии для пользователя %d: %v\n", userID, err)
				} else {
					fmt.Printf("Энергия успешно пополнена для пользователя %d\n", userID)
				}
			}

			fmt.Println("Ежедневное пополнение энергии завершено.")
		}
	}()

	http.ListenAndServe(":8080", router)
}
