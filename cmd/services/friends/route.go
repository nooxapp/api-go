package friends

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) FriendRoutes(router *mux.Router) {
	router.HandleFunc("/sendfr", h.HandleFriendRequest).Methods("POST")
}

func (h *Handler) HandleFriendRequest(w http.ResponseWriter, r *http.Request) {
	println("not done")
}
