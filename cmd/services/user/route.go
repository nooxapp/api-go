package user

import (
	"fmt"
	"net/http"
	"static-api/db"
	"static-api/helpers"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserPayload struct {
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/auth/register", h.handleReg).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload
	helpers.ReadJSON(r, &payload)
	fmt.Println(payload)
}

func (h *Handler) handleReg(w http.ResponseWriter, r *http.Request) {
	var payload UserPayload
	err := helpers.ReadJSON(r, &payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	fmt.Printf("User registered: %+v", payload)

	err = db.RegUser(payload.Username, payload.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}
