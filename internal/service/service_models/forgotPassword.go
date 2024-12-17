package service_models

type ForgotPasswordRequest struct {
	Username string `json:"username" validate:"required"`
}
