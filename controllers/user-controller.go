package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/amarjeet2003/user-api-go/models"
	"github.com/amarjeet2003/user-api-go/repository"
	"github.com/gorilla/mux"
)

type UserController struct {
	userRepo *repository.UserRepository
}

func NewUserController(userRepo *repository.UserRepository) *UserController {
	return &UserController{userRepo}
}

func validateUserInput(user *models.User) error {
	if strings.TrimSpace(user.FirstName) == "" {
		return fmt.Errorf("first name cannot be blank")
	}
	if strings.TrimSpace(user.LastName) == "" {
		return fmt.Errorf("last name cannot be blank")
	}
	if strings.TrimSpace(user.Username) == "" {
		return fmt.Errorf("username cannot be blank")
	}
	if user.DOB.After(time.Now()) {
		return fmt.Errorf("DOB cannot be in the future")
	}
	if user.DOB.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return fmt.Errorf("DOB cannot be before January 1, 1900")
	}
	if len(user.FirstName) < 2 || len(user.FirstName) > 50 {
		return fmt.Errorf("first name must be between 2 and 50 characters")
	}
	if len(user.LastName) < 2 || len(user.LastName) > 50 {
		return fmt.Errorf("last name must be between 2 and 50 characters")
	}
	if len(user.Username) < 3 || len(user.Username) > 16 {
		return fmt.Errorf("username must be between 3 and 30 characters")
	}
	for _, char := range user.Username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return fmt.Errorf("username must be alphanumeric")
		}
	}
	for _, char := range user.Username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return fmt.Errorf("username must be alphanumeric")
		}
	}

	for _, char := range user.FirstName {
		if !unicode.IsLetter(char) {
			return fmt.Errorf("first name must contain only letters")
		}
	}

	for _, char := range user.LastName {
		if !unicode.IsLetter(char) {
			return fmt.Errorf("last name must contain only letters")
		}
	}
	return nil
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reqBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	dobStr, ok := reqBody["dob"].(string)
	if !ok {
		log.Println("DOB is not a valid string")
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		log.Println("Error parsing DOB:", err)
		http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
		return
	}

	user := models.User{
		FirstName: strings.TrimSpace(reqBody["first_name"].(string)),
		LastName:  strings.TrimSpace(reqBody["last_name"].(string)),
		Username:  strings.TrimSpace(reqBody["username"].(string)),
		DOB:       dob,
	}

	if err := validateUserInput(&user); err != nil {
		log.Println("Validation error:", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	err = uc.userRepo.CreateUser(&user)
	if err != nil {
		if err.Error() == "username already exists" {
			http.Error(w, `{"error": "Username already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting ID to integer:", err)
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var reqBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}

	dobStr, ok := reqBody["dob"].(string)
	if !ok {
		log.Println("DOB is not a valid string")
		http.Error(w, `{"error": "Bad request"}`, http.StatusBadRequest)
		return
	}
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		log.Println("Error parsing DOB:", err)
		http.Error(w, `{"error": "Invalid date format"}`, http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:        id,
		FirstName: strings.TrimSpace(reqBody["first_name"].(string)),
		LastName:  strings.TrimSpace(reqBody["last_name"].(string)),
		Username:  strings.TrimSpace(reqBody["username"].(string)),
		DOB:       dob,
	}

	if err := validateUserInput(&user); err != nil {
		log.Println("Validation error:", err)
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	err = uc.userRepo.UpdateUser(&user)
	if err != nil {
		if err.Error() == "username already exists" {
			http.Error(w, `{"error": "Username already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")

	users, err := uc.userRepo.SearchUsers(name)
	if err != nil {
		log.Println("Error searching users:", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
