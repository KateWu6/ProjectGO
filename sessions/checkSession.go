package sessions

import (
	"Project/bd"
	"net/http"

	_ "github.com/lib/pq"
)

func СheckSession(next http.HandlerFunc) http.HandlerFunc {

	db, _ := bd.Connect()

	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" { // Нет куки или пустая кукушка
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		sessionID := cookie.Value

		var exists bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM sessions WHERE session_id=$1)`, sessionID).Scan(&exists)
		if err != nil {
			http.Error(w, "Ошибка при проверке сессии", http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	}
}
