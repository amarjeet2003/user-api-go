package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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
		FirstName: reqBody["first_name"].(string),
		LastName:  reqBody["last_name"].(string),
		Username:  reqBody["username"].(string),
		DOB:       dob,
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
		FirstName: reqBody["first_name"].(string),
		LastName:  reqBody["last_name"].(string),
		Username:  reqBody["username"].(string),
		DOB:       dob,
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
