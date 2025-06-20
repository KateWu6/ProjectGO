module Project

go 1.23.4

require (
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.4.0
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.39.0
)

require github.com/gorilla/securecookie v1.1.2 // indirect

replace Project/func => ./func

replace Project/handler => ./handler

replace Project/sessions => ./sessions
