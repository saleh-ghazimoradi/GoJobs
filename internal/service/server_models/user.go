package server_models

import "time"

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Email          string    `json:"email"`
	CreateAt       time.Time `json:"create_at"`
	UpdateAt       time.Time `json:"update_at"`
	IsAdmin        bool      `json:"is_admin"`
	ProfilePicture string    `json:"profile_picture"`
}
