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

// loginHandler handles user login.
// @Summary User Login
// @Description Login an existing user with username and password.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body service_models.LoginAuthPayload true "User login details"
// @Success 200 {object} service_models.User "User information with the authentication token"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// registerHandler handles the registration of a new user.
// @Summary User Registration
// @Description Register a new user with a username, password, and email.
// @Tags Users
// @Accept json
// @Produce json
// @Param body body service_models.RegisterAuthPayload true "User registration details"
// @Success 200 {object} service_models.User "Registered user details"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
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

// ForgotPasswordHandler handles forgotten passwords.
// @Summary Forgot Password
// @Description Retrieve a user's forgotten password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body service_models.ForgotPasswordRequest true "Forgot password request details"
// @Success 200 {string} string "New password"
// @Failure 400 {object} ErrorResponse "Invalid request parameters"
// @Failure 500 {object} ErrorResponse "Internal server error"
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
