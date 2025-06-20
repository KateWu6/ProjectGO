package my_check

import (
	"Project/internal/database"
)

func CheckUsername(username string) (bool, error) {
	db, err := database.Connect()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE user_name = $1", username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
