package main

type Task struct {
	ID_task          int16  `json:"id_task"`
	ID_user          int16  `json:"id_user"`
	Task_name        string `json:"task_name"`
	Task_description string `json:"task_description"`
	Time             string `json:"time"`
	Date             string `json:"date"`
	Done             bool   `json:"done"`
}

type User struct {
	ID_user     uint16 `json:"id_user"`
	Name        string `json:"user_name"`
	Password    string `json:"password"`
	Achievement uint16 `json:"achievement"`
	Energy      int16  `json:"energy"`
	Lvl         uint16 `json:"lvl"`
}

type Achievement struct {
	ID_acievement           int    `json:"id_acievement"`
	Achievement_name        string `json:"acievement_name"`
	Achievement_description string `json:"acievement_description"`
	ID_user                 int    `json:"id_user"`
	Done                    bool   `json:"done"`
}
