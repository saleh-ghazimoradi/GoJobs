package service_models

import "time"

type User struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Email          string    `json:"email"`
	CreateAt       time.Time `json:"create_at"`
	UpdateAt       time.Time `json:"update_at"`
	IsAdmin        bool      `json:"is_admin"`
	ProfilePicture *string   `json:"profile_picture"`
}

type UserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=3,max=72"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

type UpdateUserPayload struct {
	Username string `json:"username" validate:"max=100"`
	Password string `json:"password"`
	Email    string `json:"email" validate:"email,max=255"`
}
