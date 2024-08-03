package friendrequesthelper

import (
	"database/sql"
	"fmt"
	"net/http"
	"static-api/db"
	"static-api/helpers/auth"
)

func SendFriendRequest(r *http.Request, receiverUsername string) error {
	conn := db.DB

	//not the best way to get the ID
	claims, err := auth.GetSession(r)
	if err != nil {
		return fmt.Errorf("unauthorized: %w", err)
	}
	senderID := claims.ID
	var receiverID int
	err = conn.QueryRow("SELECT id FROM users WHERE username = $1", receiverUsername).Scan(&receiverID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("receiver does not exist")
		}
		return fmt.Errorf("error checking receiver existence: %w", err)
	}
	_, err = conn.Exec("INSERT INTO user_friend_requests (sender_id, receiver_id) VALUES ($1, $2)", senderID, receiverID)
	if err != nil {
		return fmt.Errorf("error sending friend request: %w", err)
	}

	return nil
}
