package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
	"github.com/saleh-ghazimoradi/GoJobs/utils"
	"path/filepath"
)

type User interface {
	GetUserById(ctx context.Context, id int64) (*service_models.User, error)
	UpdateUserProfile(ctx context.Context, id int64, username, email string) (*service_models.User, error)
	UpdateUserProfilePicture(ctx context.Context, id int64, picture string) error
	GetAllUsers(ctx context.Context) ([]*service_models.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ChangePassword(ctx context.Context, id int64, currentPassword, newPassword string) error
	GetWithTXT(tx *sql.Tx) User
}

type userService struct {
	userRepo repository.User
	tx       *sql.Tx
}

func (u *userService) GetUserById(ctx context.Context, id int64) (*service_models.User, error) {
	if u.tx != nil {
		return u.userRepo.GetWithTXT(u.tx).GetUserById(ctx, id)
	}
	return u.userRepo.GetUserById(ctx, id)
}

func (u *userService) UpdateUserProfile(ctx context.Context, id int64, username, email string) (*service_models.User, error) {
	user := &service_models.User{ID: id, Username: username, Email: email}
	return u.userRepo.UpdateUserProfile(ctx, user)
}

func (u *userService) UpdateUserProfilePicture(ctx context.Context, id int64, picture string) error {
	return u.userRepo.UpdateUserProfilePicture(ctx, id, picture)
}

func (u *userService) GetAllUsers(ctx context.Context) ([]*service_models.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *userService) DeleteUser(ctx context.Context, id int64) error {
	profilePicture, err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrRecordNotFound
		}
		return fmt.Errorf("delete user: %w", err)
	}

	if profilePicture != "" {
		filePath := filepath.Join(config.AppConfig.UploadDIR.Upload, profilePicture)
		err = utils.DeleteFileExist(filePath)
		if err != nil {
			return fmt.Errorf("delete user: %w", err)
		}
	}
	return nil
}

func (u *userService) ChangePassword(ctx context.Context, id int64, currentPassword, newPassword string) error {
	return u.userRepo.ChangePassword(ctx, id, currentPassword, newPassword)
}

func (u *userService) GetWithTXT(tx *sql.Tx) User {
	return &userService{
		userRepo: u.userRepo.GetWithTXT(tx),
		tx:       tx,
	}
}

func NewUserService(userRepo repository.User) User {
	return &userService{
		userRepo: userRepo,
	}
}
