package handler

import (
	"html/template"
	"net/http"
)

func Home_page(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html_files/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, "")
}
