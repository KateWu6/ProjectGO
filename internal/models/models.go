package models

import "time"

type Task struct {
	ID_task          int16     `json:"id_task"`
	ID_user          int16     `json:"id_user"`
	Task_name        string    `json:"task_name"`
	Task_description string    `json:"task_description"`
	Time             time.Time `json:"time"`
	Date             time.Time `json:"date"`
	Done             bool      `json:"done"`
	Exp              int       `join:"exp"`
	Energy           int       `join:"energy"`
}

type User struct {
	ID_user     uint16 `json:"id_user"`
	Name        string `json:"user_name"`
	Password    []byte `json:"password"`
	Achievement uint16 `json:"achievement"`
	Energy      int16  `json:"energy"`
	Lvl         uint16 `json:"lvl"`
	Exp         uint16 `json:"exp"`
}

type Achievement struct {
	ID_acievement           int    `json:"id_acievement"`
	Achievement_name        string `json:"acievement_name"`
	Achievement_description string `json:"acievement_description"`
	ID_user                 int    `json:"id_user"`
	Done                    bool   `json:"done"`
}

type AchievementPageData struct {
	Title          string
	IsProfile      bool
	IsTasks        bool
	IsAchievements bool
	Level          uint16
	Exp            uint16
	ExpToNextLevel int
	Energy         int16
	MaxEnergy      int
	EnergyUser     int16
	Achievements   []Achievement
	ActiveUser     string
}

type PageData struct {
	Title          string
	IsProfile      bool
	IsTasks        bool
	IsAchievements bool
	Level          uint16
	Exp            uint16
	ExpToNextLevel int
	Energy         int16
	MaxEnergy      int
	EnergyUser     int16
	User           User
}

type DayGroup struct {
	Date  string
	Tasks []Task
}

type TasksPageData struct {
	Title          string
	IsTasks        bool
	IsProfile      bool
	IsAchievements bool
	Level          uint16
	Exp            uint16
	ExpToNextLevel int
	Energy         int16
	MaxEnergy      int
	EnergyUser     int16
	GroupedTasks   []DayGroup
	ActiveUser     string
}

type CreateTaskPageData struct {
	Title          string
	IsProfile      bool
	IsTasks        bool
	IsAchievements bool
	Level          uint16
	Exp            uint16
	ExpToNextLevel int
	Energy         int16
	MaxEnergy      int
	EnergyUser     int16
	ActiveUser     string
}
