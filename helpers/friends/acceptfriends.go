package friends

import (
	"database/sql"
	"fmt"
	"static-api/db"
)

func AcceptFriendRequest(username string, currentUserID int) error {
	//this is probably a shitty way to accept friend requests
	conn := db.DB
	var senderID int
	err := conn.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user %s not found", username)
		}
		return err
	}
	var status string
	err = conn.QueryRow("SELECT status FROM user_friend_requests WHERE sender_id = $1 AND receiver_id = $2", senderID, currentUserID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no friend request found from user %s", username)
		}
		return err
	}
	if status != "pending" {
		return fmt.Errorf("friend request status is not pending")
	}
	_, err = conn.Exec("UPDATE user_friend_requests SET status = 'accepted' WHERE sender_id = $1 AND receiver_id = $2", senderID, currentUserID)
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)", senderID, currentUserID)
	if err != nil {
		return err
	}
	return nil
}
