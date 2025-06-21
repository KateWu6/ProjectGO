package services

import (
	"database/sql"
	"fmt"
)

// функция для добавления/убавления опыта/энергии
func UpdateUserExpAndEnergy(db *sql.DB, userId uint16, taskExp int, taskEnergy int) error {
	// 1. Получаем текущие значения пользователя
	var currentExp int
	var currentEnergy int
	var currentLvl uint16
	err := db.QueryRow("SELECT exp, energy, lvl FROM users WHERE id_user = $1", userId).Scan(&currentExp, &currentEnergy, &currentLvl)
	if err != nil {
		return fmt.Errorf("ошибка получения данных пользователя: %w", err)
	}

	// 2. Обновляем опыт и энергию
	newExp := currentExp + taskExp
	newEnergy := currentEnergy - taskEnergy

	// Убедимся, что энергия не ушла в минус - в будущем поменять
	if newEnergy < 0 {
		newEnergy = 0
	}

	// 3. Проверяем, нужно ли повысить уровень
	newLvl := currentLvl
	expToNextLevel := GetExpForLevel(currentLvl + 1) // Опыт для следующего уровня

	for newExp >= expToNextLevel && currentLvl < 1000 { //Максимальный уровень = 1000
		newExp -= expToNextLevel
		newLvl++
		newEnergy += 10 // Даем 10 энергии за повышение уровня
		expToNextLevel = GetExpForLevel(newLvl + 1)
	}

	// 4. Обновляем пользователя в базе данных
	_, err = db.Exec("UPDATE users SET exp = $1, energy = $2, lvl = $3 WHERE id_user = $4", newExp, newEnergy, newLvl, userId)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя в базе данных: %w", err)
	}

	return nil
}
