package friends

import (
	"net/http"
	"static-api/helpers"
	"static-api/helpers/types"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/sendfr", h.HandleFriendRequest).Methods("POST")
}

func (h *Handler) HandleFriendRequest(w http.ResponseWriter, r *http.Request) {
	var payload types.FriendRequestPayload
	//check for username in request body
	helpers.ReadJSON(r, &payload)
	//check if the payload contains "username"
	if payload.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Friend request sent successfully to " + payload.Username))
}
