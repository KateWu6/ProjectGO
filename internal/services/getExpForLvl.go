package services

func GetExpForLevel(level uint16) int {
	if level <= 10 {
		return 10
	} else if level > 10 && level <= 100 {
		return 100
	} else {
		return 1000
	}
}
