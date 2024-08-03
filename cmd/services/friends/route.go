package friends

import (
	"net/http"
	"static-api/helpers"
	"static-api/helpers/auth"
	"static-api/helpers/friends"
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
	router.HandleFunc("/acceptfr/{username}", h.HandleAcceptFriendRequest).Methods("POST")
}

func (h *Handler) HandleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	_, err := auth.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// TODO: implement the logic to accept a friend request
	vars := mux.Vars(r)
	username := vars["username"]
	friends.AcceptFriendRequest(username)
}

func (h *Handler) HandleFriendRequest(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.GetSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var payload types.FriendRequestPayload
	err = helpers.ReadJSON(r, &payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	err = friends.SendFriendRequest(r, payload.Username, claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Friend request sent successfully to " + payload.Username,
	}
	helpers.WriteJSON(w, response)
}
