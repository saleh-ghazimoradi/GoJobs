package service_models

import "time"

type Job struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	Company     string    `json:"company" validate:"required"`
	Salary      string    `json:"salary" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int64     `json:"user_id"`
}

type UpdateJobPayload struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	Company     string    `json:"company" validate:"required"`
	Salary      string    `json:"salary" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int64     `json:"user_id"`
}
