package data

import "github.com/dgrijalva/jwt-go"

type VideoEvent struct {
	VideoURL string `json:"video_url"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
}

// User struct represents the user entity
type User struct {
	ID       int    `json:"id"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	ClientID string `json:"client_id"`
	jwt.StandardClaims
}
