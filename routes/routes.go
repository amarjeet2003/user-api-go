package routes

import (
	"github.com/amarjeet2003/user-api-go/repository"
	"github.com/gorilla/mux"

	"github.com/amarjeet2003/user-api-go/controllers"
)

func SetupUserRoutes(router *mux.Router, userRepo *repository.UserRepository) {
	userController := controllers.NewUserController(userRepo)

	router.HandleFunc("/users/create", userController.CreateUser).Methods("POST")
	router.HandleFunc("/users/update/{id}", userController.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/search", userController.SearchUsers).Methods("GET")
}
