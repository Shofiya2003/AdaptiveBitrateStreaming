package data

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

type Video struct {
	VideoID    string    `json:"video_id" db:"video_id"`
	ClientID   string    `json:"client_id" db:"client_id"`
	UploadTime time.Time `json:"upload_time" db:"upload_time"`
	Status     string    `json:"status" db:"status"`
	FileKey    string    `json:"file_key" db:"file_key"`
	Bucket     string    `json:"bucket" db:"bucket"`
	Strategy   string    `json:"strategy" db:"strategy"`
	Url        string    `json:"url"`
}

type Claims struct {
	ClientID string `json:"client_id"`
	jwt.StandardClaims
}
