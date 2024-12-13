package gateway

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
	"net/http"
	"time"
)

type user struct {
	userService service.User
}

func (u *user) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
	}

	us, err := u.userService.GetUserById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			notFoundResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
		return
	}

	if err = jsonResponse(w, http.StatusOK, us); err != nil {
		internalServerError(w, r, err)
	}
}

func (u *user) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
	}

	var updateUser service_models.UpdateUserPayload

	if err = readJSON(w, r, &updateUser); err != nil {
		badRequestResponse(w, r, err)
	}

	if err = Validate.Struct(updateUser); err != nil {
		badRequestResponse(w, r, err)
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized to update this user profile"))
		return
	}

	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok {
		unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized to update this user profile"))
		return
	}

	if !isAdmin && userID != id {
		unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized to update this user profile"))
	}

	updateUse, err := u.userService.UpdateUserProfile(ctx, id, updateUser.Username, updateUser.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			notFoundResponse(w, r, err)
		default:
			internalServerError(w, r, err)
		}
	}
	if err = jsonResponse(w, http.StatusOK, updateUse); err != nil {
		internalServerError(w, r, err)
	}
}

func NewUserHandler(userService service.User) *user {
	return &user{
		userService: userService,
	}
}
