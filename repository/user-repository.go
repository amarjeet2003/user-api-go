package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/amarjeet2003/user-api-go/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) CreateUser(user *models.User) error {
	isUnique, err := ur.IsUsernameUnique(user.Username, 0)
	if err != nil {
		log.Println("Error checking username uniqueness:", err)
		return err
	}
	if !isUnique {
		return fmt.Errorf("username already exists")
	}

	query := "INSERT INTO users (first_name, last_name, username, dob) VALUES (?, ?, ?, ?)"
	result, err := ur.db.Exec(query, user.FirstName, user.LastName, user.Username, user.DOB.Format("2006-01-02"))
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return err
	}
	user.ID = int(lastInsertID)
	return nil
}

func (ur *UserRepository) SearchUsers(name string) ([]*models.User, error) {
	var query string
	var rows *sql.Rows
	var err error

	if name != "" {
		query = "SELECT id, first_name, last_name, username, dob FROM users WHERE first_name LIKE ? OR last_name LIKE ?"
		rows, err = ur.db.Query(query, "%"+name+"%", "%"+name+"%")
	} else {
		query = "SELECT id, first_name, last_name, username, dob FROM users ORDER BY dob"
		rows, err = ur.db.Query(query)
	}

	if err != nil {
		log.Println("Error searching users:", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var dob string

		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &dob); err != nil {
			log.Println("Error scanning user row:", err)
			return nil, err
		}

		user.DOB, err = time.Parse("2006-01-02", dob)
		if err != nil {
			log.Println("Error parsing date:", err)
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
	var dob string

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &dob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user.DOB, err = time.Parse("2006-01-02", dob)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) IsUsernameUnique(username string, userID int) (bool, error) {
	var id int
	query := "SELECT id FROM users WHERE username = ? AND id != ?"
	err := ur.db.QueryRow(query, username, userID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (ur *UserRepository) UpdateUser(user *models.User) error {
	isUnique, err := ur.IsUsernameUnique(user.Username, user.ID)
	if err != nil {
		log.Println("Error checking username uniqueness:", err)
		return err
	}
	if !isUnique {
		return fmt.Errorf("username already exists")
	}

	query := "UPDATE users SET first_name=?, last_name=?, username=?, dob=? WHERE id=?"
	_, err = ur.db.Exec(query, user.FirstName, user.LastName, user.Username, user.DOB.Format("2006-01-02"), user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return err
	}
	return nil
}
