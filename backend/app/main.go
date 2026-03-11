package main

import (
	"log"
	"net/http"
	"quiz-irr/internals/database"
	"quiz-irr/internals/handlers"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/middlewares"
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
	testsRepo := storage.NewTestsRepo(db)
	questionsRepo := storage.NewQuestionsRepo(db)
	optionsRepo := storage.NewOptionsRepo(db)
	rawRepo := storage.NewRawRepo(db)

	usersService := services.NewUsersService(usersRepo)
	authService := services.NewAuthServce("FDGSSDFGSDFG")
	testsService := services.NewTestService(testsRepo)
	questionsService := services.NewQuestionService(questionsRepo)
	optionsService := services.NewOptionService(optionsRepo)
	rawService := services.NewRawService(rawRepo)

	testCases := usecases.NewTestsCases(testsService, optionsService, questionsService)
	userCases := usecases.NewUsersCases(usersService, authService)
	examCases := usecases.NewExamCases(testsService, rawService, optionsService, questionsService)

	testHandler := handlers.NewTestHandlers(testCases, userCases)
	authHandlers := handlers.NewAuthHandlers(userCases)
	userHandlers := handlers.NewUserHandlers(userCases)
	examHandlers := handlers.NewExamHandlers(examCases)

	authMiddlware := middlewares.NewAuthMiddleware(authService)

	r := router.NewRouter(
		authHandlers,
		userHandlers,
		testHandler,
		examHandlers,
		authMiddlware,
	)
	root, err := userCases.BootstrapCreate(dto.CreateUserRequest{
		FullName: "Абрамов Вячеслав Александрович",
		Password: "changeme",
		Email:    "vyachik005@gmail.com",
	})
	if err != nil {
		log.Println("Failed admin creating: " + err.Error())
	} else {
		log.Println("Root Name: ", root.FullName)
		log.Println("Root Email: ", root.Email)
	}

	log.Println("Server Started...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
