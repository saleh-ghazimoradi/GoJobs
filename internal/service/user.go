package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
)

type User interface {
	GetUserById(ctx context.Context, id int64) (*service_models.User, error)
	GetWithTXT(tx *sql.Tx) User
}

type userService struct {
	userRepo repository.User
}

func (u *userService) GetUserById(ctx context.Context, id int64) (*service_models.User, error) {
	return u.userRepo.GetUserById(ctx, id)
}

func (u *userService) GetWithTXT(tx *sql.Tx) User {
	return &userService{
		userRepo: u.userRepo.GetWithTXT(tx)}
}

func NewUserService(userRepo repository.User) User {
	return &userService{
		userRepo: userRepo,
	}
}