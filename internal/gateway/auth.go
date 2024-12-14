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
