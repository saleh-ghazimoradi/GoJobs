package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service/service_models"
)

type User interface {
	CreateUser(ctx context.Context, user *service_models.User) error
	GetUserById(ctx context.Context, id int64) (*service_models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*service_models.User, error)
	UpdateUserProfile(ctx context.Context, user *service_models.User) (*service_models.User, error)
	UpdateUserProfilePicture(ctx context.Context, id int64, picture string) error
	GetAllUsers(ctx context.Context) ([]*service_models.User, error)
	UpdateUserPassword(ctx context.Context, user *service_models.User) error
	DeleteUser(ctx context.Context, id int64) (string, error)
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

func (u *userRepository) UpdateUserProfilePicture(ctx context.Context, id int64, picture string) error {
	query := `UPDATE users SET profile_picture = $1 WHERE id = $2`
	_, err := u.dbWrite.ExecContext(ctx, query, picture, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetAllUsers(ctx context.Context) ([]*service_models.User, error) {
	var users []*service_models.User
	query := `SELECT id, username, password, email, created_at, updated_at, is_admin, profile_picture FROM users`
	rows, err := u.dbRead.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user service_models.User
		var profilePicture sql.NullString
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreateAt, &user.UpdateAt, &user.IsAdmin, &profilePicture)
		if err != nil {
			return nil, err
		}
		if profilePicture.Valid {
			user.ProfilePicture = &profilePicture.String
		} else {
			user.ProfilePicture = nil
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *userRepository) UpdateUserPassword(ctx context.Context, user *service_models.User) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := u.dbWrite.ExecContext(ctx, query, user.Password, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id int64) (string, error) {
	query := `DELETE FROM users WHERE id = $1`
	result, err := u.dbWrite.ExecContext(ctx, query, id)
	if err != nil {
		return "", fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("error getting rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return "", ErrRecordNotFound
	}

	var profilePicture sql.NullString
	query = `SELECT profile_picture FROM users WHERE id = $1`
	err = u.dbRead.QueryRowContext(ctx, query, id).Scan(&profilePicture)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("error retrieving profile picture: %v", err)
	}

	return profilePicture.String, nil
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
