package friends

import (
	"net/http"
	"static-api/helpers"
	"static-api/helpers/auth"
	friendrequesthelper "static-api/helpers/friends"
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

	err = friendrequesthelper.SendFriendRequest(r, payload.Username, claims.ID)
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
