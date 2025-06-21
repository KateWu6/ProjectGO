package main

import (
	"Project/handler"
	"Project/internal/database"
	"Project/internal/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter() //создание нового роутера

	//регистрация обработчиков маршрутов
	router.HandleFunc("/", handler.ProfilePage).Methods("GET")

	router.HandleFunc("/login", handler.LoginHandler)
	router.HandleFunc("/register", handler.RegisterHandler)
	router.HandleFunc("/logout", handler.LogoutHandler).Methods("POST")

	router.HandleFunc("/tasks/", handler.GetUserTasks)
	router.HandleFunc("/achievements/", handler.GetUserAchievement)
	//маршрут для создания новой задачи
	router.HandleFunc("/tasks/create_task/", handler.CreateTask)
	//маршрут обрабатывающий завершение задачи с идентификатором {id}
	router.HandleFunc("/tasks/{id}/complete", handler.CompleteTaskHandler).Methods("POST")

	//запуск фоновой функции (востановление энергии каждый день)
	go func() {
		for {
			// Получение списка всех user_id из базы данных
			db, err := database.Connect()
			if err != nil {
				fmt.Println("Ошибка подключения к базе данных:", err)
				time.Sleep(24 * time.Hour) // Повторная попытка через 24 часа
				continue
			}
			//получение списка всех id из таблицы users
			rows, err := db.Query("SELECT id_user FROM users") //итерация по результатам запроса
			if err != nil {
				fmt.Println("Ошибка получения списка пользователей:", err)
				db.Close()
				time.Sleep(24 * time.Hour) // Повторная попытка через 24 часа
				continue
			}
			defer rows.Close()

			var userIDs []uint16 //массив идентификаторов пользователей
			//извлечение результатов запроса в цикле
			for rows.Next() {
				var userID uint16
				if err := rows.Scan(&userID); err != nil {
					fmt.Println("Ошибка сканирования ID пользователя:", err)
					continue
				}
				userIDs = append(userIDs, userID)
			}
			db.Close()

			now := time.Now() //текущая дата
			//определение ближайшего к полуночи времени
			midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			nextMidnight := midnight.Add(24 * time.Hour)   //Время следующей полуночи
			durationUntilMidnight := nextMidnight.Sub(now) //рассчет времени до следующей полуночи

			time.Sleep(durationUntilMidnight) //"засыпание" программы на указанное время
			//восстановление энергии пользователям
			for _, userID := range userIDs {
				if err := services.ReplenishEnergy(userID); err != nil {
					fmt.Printf("Ошибка пополнения энергии для пользователя %d: %v\n", userID, err)
				} else {
					fmt.Printf("Энергия успешно пополнена для пользователя %d\n", userID)
				}
			}

			fmt.Println("Ежедневное пополнение энергии завершено.")
		}
	}()

	http.ListenAndServe(":8080", router) //вызов функции запуска сервера
}
