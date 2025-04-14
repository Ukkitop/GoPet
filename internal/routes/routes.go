package routes

import (
	"github.com/gorilla/mux"
	"realAPI/internal/database"
	"realAPI/internal/handlers"
	"realAPI/internal/middleware"
)

func Routes(r *mux.Router, client *database.Server) {

	r.StrictSlash(true)
	r.HandleFunc("/user/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/user/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/user/logout", middleware.Authorization(handlers.LogoutHandler)).Methods("GET")
	r.HandleFunc("/mongo/healthcheck", client.HealthCheck).Methods("GET")
	r.HandleFunc("/mongo/addword", client.AddWord).Methods("POST")
	r.HandleFunc("/mongo/updateword", client.UpdateWord).Methods("POST")
	r.HandleFunc("/ws", handlers.WsEndpoint)
}
