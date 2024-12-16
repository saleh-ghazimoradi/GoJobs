package gateway

import (
	"context"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
	"net/http"
	"time"
)

type authenticate struct {
	authService service.Authenticate
}

// loginHandler handles user login and returns a token.
// @Summary User login
// @Description Authenticates a user and returns a JWT token if the credentials are valid.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param LoginAuthPayload body service_models.LoginAuthPayload true "Login credentials"
// @Success 200 {string} string "JWT token"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/login [post]
func (a *authenticate) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var loginAuthPayload service_models.LoginAuthPayload
	if err := readJSON(w, r, &loginAuthPayload); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(loginAuthPayload); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	token, err := a.authService.LoginUser(ctx, loginAuthPayload.Username, loginAuthPayload.Password)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, token); err != nil {
		internalServerError(w, r, err)
	}

}

// registerHandler handles user registration by creating a new user.
// @Summary User registration
// @Description Registers a new user with the provided username, password, and email.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param RegisterAuthPayload body service_models.RegisterAuthPayload true "User registration credentials"
// @Success 201 {object} service_models.User "User created"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/register [post]
func (a *authenticate) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var registerAuthPayload service_models.RegisterAuthPayload
	if err := readJSON(w, r, &registerAuthPayload); err != nil {
		badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(registerAuthPayload); err != nil {
		badRequestResponse(w, r, err)
	}

	us := &service_models.User{
		Username: registerAuthPayload.Username,
		Password: registerAuthPayload.Password,
		Email:    registerAuthPayload.Email,
	}

	if err := a.authService.RegisterUser(ctx, us); err != nil {
		internalServerError(w, r, err)
	}

	if err := jsonResponse(w, http.StatusCreated, us); err != nil {
		internalServerError(w, r, err)
	}
}

// ForgotPasswordHandler handles the password reset request by sending a password to the user.
// @Summary Password reset request
// @Description Requests a password reset for the provided username and returns a password if successful.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param ForgotPasswordRequest body service_models.ForgotPasswordRequest true "User's username for password reset"
// @Success 200 {string} string "Password reset successful"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/forgotpassword [post]
func (a *authenticate) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var passReq service_models.ForgotPasswordRequest
	if err := readJSON(w, r, &passReq); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(passReq); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	password, err := a.authService.ForgotPassword(ctx, passReq.Username)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	if err = jsonResponse(w, http.StatusOK, password); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func NewAuthenticateHandler(authService service.Authenticate) *authenticate {
	return &authenticate{authService: authService}
}
