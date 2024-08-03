package friends

import (
	"database/sql"
	"fmt"
	"static-api/db"
)

func AcceptFriendRequest(username string) error {
	conn := db.DB
	println(username)

	// Connect to db and update the enum
	var status string
	err := conn.QueryRow("SELECT status FROM user_friend_requests WHERE username = $1", username).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("bad request: no friend request found for user %s", username)
		}
		return err
	}

	return nil
}
