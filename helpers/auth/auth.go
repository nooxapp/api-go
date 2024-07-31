package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"static-api/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type User struct {
	Username string
	Email    string
	Password string
}

func RegUser(Username, Email, Password string) error {
	conn := db.DbConnect()
	defer conn.Close()

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

func AuthUser(Email, Password string) (*User, error) {
	conn := db.DbConnect()
	defer conn.Close()

	var user User
	query := `SELECT username, email, password FROM users WHERE email = $1`
	err := conn.QueryRow(query, Email).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &user, nil
}

func GenerateJWT(Email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
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
