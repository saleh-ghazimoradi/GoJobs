package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
)

type User interface {
	CreateUser(ctx context.Context, user *service_models.User) error
	GetUserById(ctx context.Context, id int64) (*service_models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*service_models.User, error)
	UpdateUserProfile(ctx context.Context, user *service_models.User) (*service_models.User, error)
	GetWithTXT(tx *sql.Tx) User
}

type userRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
	tx      *sql.Tx
}

func (u *userRepository) CreateUser(ctx context.Context, user *service_models.User) error {
	query := `INSERT INTO users(username,password,email) VALUES ($1,$2,$3) RETURNING id, created_at`

	err := u.dbWrite.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&user.ID, &user.CreateAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmails
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsernames
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetUserById(ctx context.Context, id int64) (*service_models.User, error) {
	var user service_models.User
	var profilePicture sql.NullString
	query := `SELECT id, username, password, email, created_at, updated_at, is_admin, profile_picture FROM users WHERE id = $1`

	err := u.dbRead.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreateAt, &user.UpdateAt, &user.IsAdmin, &profilePicture)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	if profilePicture.Valid {
		user.ProfilePicture = &profilePicture.String
	} else {
		user.ProfilePicture = nil
	}
	return &user, err
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*service_models.User, error) {
	var user service_models.User
	query := `SELECT id, username, password, email, created_at, updated_at, is_admin, profile_picture FROM users WHERE username = $1`

	err := u.dbRead.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreateAt, &user.UpdateAt, &user.IsAdmin, &user.ProfilePicture)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, err
}

func (u *userRepository) UpdateUserProfile(ctx context.Context, user *service_models.User) (*service_models.User, error) {
	query := `UPDATE users SET username = $1, email = $2 WHERE id = $3`
	_, err := u.dbWrite.ExecContext(ctx, query, user.Username, user.Email, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userRepository) GetWithTXT(tx *sql.Tx) User {
	return &userRepository{
		dbWrite: u.dbWrite,
		dbRead:  u.dbRead,
		tx:      tx,
	}
}

func NewUserRepository(dbWrite *sql.DB, dbRead *sql.DB) User {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
