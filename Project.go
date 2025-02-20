package main

import (
	"fmt"
	"net/http"
)

type Task struct {
	task_name string
	time      string
	date      string
}

type User struct {
	name        string
	password    string
	achievement uint16
	energy      int16
	lvl         uint16
	tasks       []Task
}

func (u User) getAllInfo() string {
	return fmt.Sprintf("You name is: %s. Your password is %s and you "+
		"have %d achievement", u.name, u.password, u.achievement)
}

func (u *User) setNewName(newName string) {
	u.name = newName
}

func home_page(w http.ResponseWriter, r *http.Request) {
	Bob := User{"Bob", "1234", 100, 10, 1, []Task{{"sport", "6.00", "01.02.25"}}}
	Bob.setNewName("Alex")
	fmt.Fprintf(w, Bob.getAllInfo())
}

func contacts_page(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Contacts")
}

func handlRequest() {
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8080", nil)
}

func main() {
	//Bob := User{"Bob","1234", 0, 10,  1, []string{"not"}}
	handlRequest()
}
