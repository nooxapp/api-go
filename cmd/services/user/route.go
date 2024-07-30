package user

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", h.handleLogin).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	//Read Request Body
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	println(string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Request body received"}`))
}
