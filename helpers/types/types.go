package types

type UserPayload struct {
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type FriendRequestPayload struct {
	Username string `json:"Username"`
}
