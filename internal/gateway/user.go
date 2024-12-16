package gateway

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
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

// getUserByIdHandler retrieves a user by ID.
// @Summary Get user by ID
// @Description Fetch a user's details using their unique ID. Requires an authorization token.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} service_models.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/users/{id} [get]
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

// UpdateUserProfileHandler updates the profile of a user by ID.
// @Summary Update user profile
// @Description Update the user profile (username and email). Requires authorization token and admin check.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param updateUser body service_models.UpdateUserPayload true "User Profile Information"
// @Success 200 {object} service_models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/users/{id} [put]
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

// UpdateUserProfilePictureHandler updates the profile picture of a user by ID.
// @Summary Update user profile picture
// @Description Update the user's profile picture. Requires authorization token and admin check.
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param profile_picture formData file true "Profile Picture File"
// @Success 200 {string} string "Profile picture updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/users/{id}/picture [put]
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

// GetAllUsersHandler handles the retrieval of all users.
// @Summary Retrieve all users
// @Description Get a list of all users. Requires an authorization token.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} service_models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/users [get]
func (u *user) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok {
		unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized to get all users"))
		return
	}

	if !isAdmin {
		unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized to get all users"))
		return
	}

	users, err := u.userService.GetAllUsers(ctx)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	if err = jsonResponse(w, http.StatusOK, users); err != nil {
		internalServerError(w, r, err)
		return
	}
}

// DeleteUserHandler deletes a user by ID.
// @Summary Delete user
// @Description Deletes a user by ID. Only an admin user can delete another user. You cannot delete yourself.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/users/{id} [delete]
func (u *user) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	isAdmin, ok := r.Context().Value("isAdmin").(bool)
	if !ok || !isAdmin {
		unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized to delete user"))
		return
	}

	id, err := readIDParam(r)
	if err != nil {
		badRequestResponse(w, r, err)
		return
	}

	currentUserID := r.Context().Value("userID").(int64)
	if currentUserID == id {
		badRequestResponse(w, r, fmt.Errorf("cannot delete yourself"))
		return
	}

	err = u.userService.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			badRequestResponse(w, r, fmt.Errorf("user not found"))
			return
		}
		internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, "user deleted"); err != nil {
		internalServerError(w, r, err)
		return
	}
}

// ChangePasswordHandler changes the user's password.
// @Summary Change password
// @Description Changes the password for the authenticated user. The user must provide their current password and the new password.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ChangePassword body service_models.ChangePassword true "Change Password Request"
// @Success 200 {string} string "Password successfully changed"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/users/{id}/changePassword [put]
func (u *user) ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var req service_models.ChangePassword
	if err := readJSON(w, r, &req); err != nil {
		badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(req); err != nil {
		badRequestResponse(w, r, err)
		return
	}

	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		badRequestResponse(w, r, fmt.Errorf("unauthorized to update this user profile"))
		return
	}

	err := u.userService.ChangePassword(ctx, userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			notFoundResponse(w, r, err)
			return
		default:
			internalServerError(w, r, err)
			return
		}
	}
	if err := jsonResponse(w, http.StatusOK, "Password successfully changed"); err != nil {
		internalServerError(w, r, err)
		return
	}
}

func NewUserHandler(userService service.User) *user {
	return &user{
		userService: userService,
	}
}
