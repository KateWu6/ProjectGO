package my_check

import (
	"Project/internal/database"
	"Project/internal/models"
	"database/sql"
)

func CheckUserExists(username string) (*models.User, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var user models.User
	err = db.QueryRow("SELECT user_name, password FROM users WHERE user_name = $1", username).Scan(&user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
