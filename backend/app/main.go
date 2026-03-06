package main

import (
	"log"
	"net/http"
	"quiz-irr/internals/database"
	"quiz-irr/internals/handlers"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/router"
	"quiz-irr/internals/services"
	"quiz-irr/internals/storage"
	"quiz-irr/internals/usecases"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed connect: " + err.Error())
	}

	usersRepo := storage.NewUsersRepo(db)
	usersService := services.NewUsersService(usersRepo)
	authService := services.NewAuthServce("FDGSSDFGSDFG")
	userCases := usecases.NewUsersCases(usersService, authService)
	userHandlers := handlers.NewUserHandlers(userCases)
	r := router.NewRouter(userHandlers)

	root, err := userCases.BootstrapCreate(dto.CreateUserRequest{
		FullName: "Абрамов Вячеслав Александрович",
		Password: "changeme",
		Email:    "vyachik005@gmail.com",
	})
	if err != nil {
		log.Println("Failed admin creating: " + err.Error())
	} else {
		log.Println("Root Name: ", root.FullName)
	}

	log.Println("Server Started...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
