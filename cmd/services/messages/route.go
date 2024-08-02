package messages

import (
	"fmt"
	"net/http"
	"static-api/helpers"
	"static-api/helpers/auth"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/msg/sendmessage", h.SendMessage).Methods("POST")
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	message := fmt.Sprintf("Hi %d", claims.ID)
	helpers.WriteJSON(w, message)
}
