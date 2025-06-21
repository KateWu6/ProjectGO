package services

import (
	"Project/internal/database"
)

const (
	MaxEnergy       = 100 // Максимальная энергия
	EnergyReplenish = 10  // Размер "баночки"
)

func ReplenishEnergy(userID uint16) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	var currentEnergy int16
	err = db.QueryRow("SELECT energy FROM users WHERE id_user = $1", userID).Scan(&currentEnergy)
	if err != nil {
		return err
	}
	//поменять потом
	if currentEnergy < EnergyReplenish {
		newEnergy := EnergyReplenish
		if newEnergy > MaxEnergy {
			newEnergy = MaxEnergy
		}

		_, err = db.Exec("UPDATE users SET energy = $1 WHERE id_user = $2", newEnergy, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
