package gateway

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func (u *user) UpdateUserProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	id, err := readIDParam(r)
	if err != nil {
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

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	file, header, err := r.FormFile("profile_picture")
	if err != nil {
		badRequestResponse(w, r, err)
	}
	defer file.Close()

	uploadDir := config.AppConfig.UploadDIR.Upload
	if err = os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		internalServerError(w, r, err)
		return
	}

	filename := fmt.Sprintf("%d-%s", id, filepath.Base(header.Filename))
	filePath := filepath.Join(uploadDir, filename)

	saveFile, err := os.Create(filePath)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	defer saveFile.Close()

	if _, err := io.Copy(saveFile, file); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := u.userService.UpdateUserProfilePicture(ctx, id, filename); err != nil {
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, "profile picture updated successfully"); err != nil {
		internalServerError(w, r, err)
	}
}

func NewUserHandler(userService service.User) *user {
	return &user{
		userService: userService,
	}
}
