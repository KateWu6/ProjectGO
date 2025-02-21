package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Task struct {
	Task_name string
	Time      string
	Date      string
}

type User struct {
	Name        string
	Password    string
	Achievement uint16
	Energy      int16
	Lvl         uint16
	Tasks       []Task
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("You name is: %s. Your password is %s and you "+
		"have %d achievement", u.Name, u.Password, u.Achievement)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
	Bob := User{"Bob", "1234", 100, 10, 1, []Task{{"sport", "6.00", "01.02.25"}}}
	//fmt.Fprintf(w, "<b>Main Text</b>")
	tmpl, _ := template.ParseFiles("templates/home_page.html")
	tmpl.Execute(w, Bob)
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contacts")
}

func handlRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8000", nil)
}

func main() {
	//Bob := User{"Bob","1234", 0, 10,  1, []string{"not"}}
	handlRequest()
}
