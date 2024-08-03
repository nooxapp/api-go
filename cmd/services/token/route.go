package token

import (
	"encoding/json"
	"net/http"
	"static-api/db"
	"static-api/helpers/auth"
	"static-api/helpers/friends"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/token", h.GetUser).Methods("POST")
}
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	conn := db.DB
	claims, err := auth.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token cookie not found", http.StatusBadRequest)
		return
	}

	userID := claims.ID
	friendsList := friends.GetFriends(userID)

	rows, err := conn.Query("SELECT id, username, email FROM users WHERE id = $1", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var id int
	var username, email string
	if rows.Next() {
		err := rows.Scan(&id, &username, &email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"Token":   tokenCookie.Value,
		"user": map[string]interface{}{
			"id":       id,
			"username": username,
			"email":    email,
		},
		"Friends": friendsList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
