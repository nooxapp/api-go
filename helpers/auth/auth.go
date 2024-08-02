package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"static-api/db"
	"static-api/helpers/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func RegUser(Username, Email, Password string) error {
	conn := db.DB

	var usernameCount int
	checkUsernameQuery := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := conn.QueryRow(checkUsernameQuery, Username).Scan(&usernameCount)
	if err != nil {
		return err
	}

	if usernameCount > 0 {
		return fmt.Errorf("username already exists")
	}
	var emailCount int
	checkEmailQuery := `SELECT COUNT(*) FROM users WHERE email = $1`
	err = conn.QueryRow(checkEmailQuery, Email).Scan(&emailCount)
	if err != nil {
		return err
	}

	if emailCount > 0 {
		return fmt.Errorf("email already exists")
	}
	insertQuery := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, execErr := conn.Exec(insertQuery, Username, Email, Password)
	if execErr != nil {
		return execErr
	}

	return nil
}

func AuthUser(Email, Password string) (*types.UserPayload, int, error) {
	conn := db.DB

	var user types.UserPayload
	var userID int
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	err := conn.QueryRow(query, Email).Scan(&userID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, fmt.Errorf("invalid credentials")
		}
		return nil, 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	if err != nil {
		return nil, 0, fmt.Errorf("invalid credentials")
	}

	return &user, userID, nil
}

func GenerateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	conn := db.DB

	insertSessionQuery := `INSERT INTO users_session (token, expires_at, user_id) VALUES ($1, $2, $3)`
	_, execErr := conn.Exec(insertSessionQuery, tokenString, expirationTime, userID)
	if execErr != nil {
		return "", execErr
	}

	return tokenString, nil
}

func GetSession(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, fmt.Errorf("no token found")
		}
		return nil, err
	}
	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
