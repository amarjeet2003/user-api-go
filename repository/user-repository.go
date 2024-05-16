package repository

import (
	"database/sql"
	"log"

	"github.com/amarjeet2003/user-api-go/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	query := "INSERT INTO users (first_name, last_name, username, dob) VALUES (?, ?, ?, ?)"
	_, err := ur.db.Exec(query, user.FirstName, user.LastName, user.Username, user.DOB)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	return nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	query := "UPDATE users SET first_name=?, last_name=?, username=?, dob=? WHERE id=?"
	_, err := ur.db.Exec(query, user.FirstName, user.LastName, user.Username, user.DOB, user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}

func (ur *UserRepository) SearchUsersByName(name string) ([]*models.User, error) {
	query := "SELECT * FROM users WHERE first_name LIKE ? OR last_name LIKE ?"
	rows, err := ur.db.Query(query, "%"+name+"%", "%"+name+"%")
	if err != nil {
		log.Println("Error searching users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.DOB); err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over user rows:", err)
		return nil, err
	}
	return users, nil
}

func (ur *UserRepository) GetUser(identifier string) (*models.User, error) {
	query := "SELECT id, first_name, last_name, username, dob FROM users WHERE username = ? OR id = ?"
	row := ur.db.QueryRow(query, identifier, identifier)
	return ur.scanUser(row)
}

func (ur *UserRepository) scanUser(row *sql.Row) (*models.User, error) {
	var user models.User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.DOB)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err
	}
	return &user, nil
}

// TODO: Update API
