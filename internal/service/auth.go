package service

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
	"github.com/saleh-ghazimoradi/GoJobs/utils"
	"golang.org/x/crypto/bcrypt"
)

type Authenticate interface {
	RegisterUser(ctx context.Context, user *service_models.User) error
	LoginUser(ctx context.Context, username, password string) (string, error)
	ForgotPassword(ctx context.Context, username string) (string, error)
	GetWithTXT(tx *sql.Tx) Authenticate
}

type authService struct {
	userRepo repository.User
}

func (a *authService) RegisterUser(ctx context.Context, user *service_models.User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return a.userRepo.CreateUser(ctx, user)
}

func (a *authService) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := a.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	return utils.GenerateToken(user.Username, user.ID, user.IsAdmin)
}

func (a *authService) GetWithTXT(tx *sql.Tx) Authenticate {
	return &authService{
		userRepo: a.userRepo.GetWithTXT(tx),
	}
}

func (a *authService) ForgotPassword(ctx context.Context, username string) (string, error) {
	user, err := a.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	generatedPassword := utils.GeneratePassword(6)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashPassword)

	if err = a.userRepo.UpdateUserPassword(ctx, user); err != nil {
		return "", err
	}
	return generatedPassword, nil
}

func NewAuthenticateService(userRepo repository.User) Authenticate {
	return &authService{
		userRepo: userRepo,
	}
}
