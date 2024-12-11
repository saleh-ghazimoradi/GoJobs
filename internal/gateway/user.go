package gateway

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
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

func NewUserHandler(userService service.User) *user {
	return &user{
		userService: userService,
	}
}
