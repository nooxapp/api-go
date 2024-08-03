package friends

import (
	"database/sql"
	"fmt"
	"net/http"
	"static-api/db"
)

func SendFriendRequest(r *http.Request, receiverUsername string, senderID int) error {
	//figured out how to get the userid from the request
	conn := db.DB
	var receiverID int
	err := conn.QueryRow("SELECT id FROM users WHERE username = $1", receiverUsername).Scan(&receiverID)
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

func GetFriends(userID int) []int {
	conn := db.DB
	rows, err := conn.Query("SELECT friend_id FROM friends WHERE user_id = $1", userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var friends []int
	for rows.Next() {
		var friendID int
		err := rows.Scan(&friendID)
		if err != nil {
			continue
		}
		friends = append(friends, friendID)
	}

	return friends
}
