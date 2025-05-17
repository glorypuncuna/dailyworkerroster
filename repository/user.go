package repository

import (
	"dailyworkerroster/model"
	"database/sql"
)

type UserRepoItf interface {
	SignUp(user *model.User) (int64, error)
	Login(identifier string) (*model.User, error)
	GetUsersByRole(role string) ([]*model.User, error)
	GetUserByID(id int64) (*model.User, error)
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepoItf {
	return &UserRepository{
		DB: db,
	}
}

// SignUp inserts a new user into the user_account table
func (r *UserRepository) SignUp(user *model.User) (int64, error) {
	query := `
        INSERT INTO user_account (name, username, email, password, role, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, NOW(), NOW())
    `
	result, err := r.DB.Exec(query, user.Name, user.Username, user.Email, user.Password, user.Role)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Login checks if a user exists with the given username/email and password
func (r *UserRepository) Login(identifier string) (*model.User, error) {
	query := `
        SELECT id, name, username, email, password, role, created_at, updated_at
        FROM user_account
        WHERE (username = ? OR email = ?)
    `
	var user model.User
	err := r.DB.QueryRow(query, identifier, identifier).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsersByRole(role string) ([]*model.User, error) {
	query := `
        SELECT id, name, username, email, password, role, created_at, updated_at
        FROM user_account
        WHERE role = ?
    `
	rows, err := r.DB.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.Name, &user.Username, &user.Email, &user.Password,
			&user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int64) (*model.User, error) {
	query := `
        SELECT id, name, username, email, password, role, created_at, updated_at
        FROM user_account
        WHERE id = ?
        LIMIT 1
    `
	var user model.User
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
