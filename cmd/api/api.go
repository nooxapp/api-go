package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"static-api/cmd/services/user"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	userService := user.NewHandler()
	userService.RegisterRoutes(subrouter)
	fmt.Println("Listening on http://localhost" + s.addr + "/api/v1/")
	return http.ListenAndServe(s.addr, router)
}
