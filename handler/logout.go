package handler

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Удаляем сессию
	session, _ := store.Get(r, "session-name")
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Перенаправляем на страницу логина или главную
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
