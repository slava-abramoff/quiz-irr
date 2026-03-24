package main

import (
	"log"
	"net/http"
	"os"
	"quiz-irr/internals/database"
	"quiz-irr/internals/handlers"
	"quiz-irr/internals/handlers/dto"
	"quiz-irr/internals/middlewares"
	"quiz-irr/internals/router"
	"quiz-irr/internals/services"
	"quiz-irr/internals/storage"
	"quiz-irr/internals/usecases"

	"github.com/joho/godotenv"
)

func main() {
	var env bool
	if err := godotenv.Load(); err != nil {
		env = false
	} else {
		env = true
	}

	var secret string
	var port string
	var root dto.CreateUserRequest

	if env {
		secret = os.Getenv("SECRET")
		port = os.Getenv("BACKEND_PORT")
		root = dto.CreateUserRequest{
			FullName: os.Getenv("ADMIN_NAME"),
			Password: os.Getenv("ADMIN_PASSWORD"),
			Email:    os.Getenv("ADMIN_EMAIL"),
		}
	} else {
		secret = "FDGSSDFGSDFG"
		port = ":8080"
		root = dto.CreateUserRequest{
			FullName: "Абрамов Вячеслав Александрович",
			Password: "changeme",
			Email:    "vyachik005@gmail.com",
		}
	}

	db, err := database.ConnectDB(env)
	if err != nil {
		log.Fatal("Failed connect: " + err.Error())
	}

	usersRepo := storage.NewUsersRepo(db)
	testsRepo := storage.NewTestsRepo(db)
	questionsRepo := storage.NewQuestionsRepo(db)
	optionsRepo := storage.NewOptionsRepo(db)
	rawRepo := storage.NewRawRepo(db)
	resultsRepo := storage.NewResultRepo(db)

	usersService := services.NewUsersService(usersRepo)
	authService := services.NewAuthServce(secret)
	testsService := services.NewTestService(testsRepo)
	questionsService := services.NewQuestionService(questionsRepo)
	optionsService := services.NewOptionService(optionsRepo)
	rawService := services.NewRawService(rawRepo)
	resultsService := services.NewResultsService(resultsRepo)
	excelService := services.NewExcelService()

	testCases := usecases.NewTestsCases(testsService, optionsService, questionsService)
	userCases := usecases.NewUsersCases(usersService, authService)
	examCases := usecases.NewExamCases(testsService, rawService, optionsService, questionsService)
	rawAnswersCases := usecases.NewRawAnswersCases(
		rawService,
		questionsService,
		optionsService,
		resultsService,
	)
	resultsCases := usecases.NewResultsCases(resultsService, excelService)

	testHandler := handlers.NewTestHandlers(testCases, userCases)
	authHandlers := handlers.NewAuthHandlers(userCases)
	userHandlers := handlers.NewUserHandlers(userCases)
	examHandlers := handlers.NewExamHandlers(examCases)
	rawAnswersHandlers := handlers.NewRawAnswersHandlers(rawAnswersCases)
	resultsHandlers := handlers.NewResultsHandlers(resultsCases)

	authMiddlware := middlewares.NewAuthMiddleware(authService)

	r := router.NewRouter(
		authHandlers,
		userHandlers,
		testHandler,
		examHandlers,
		rawAnswersHandlers,
		resultsHandlers,
		authMiddlware,
	)
	newRoot, err := userCases.BootstrapCreate(root)
	if err != nil {
		log.Println("Failed admin creating: " + err.Error())
	} else {
		log.Println("Root Name: ", newRoot.FullName)
		log.Println("Root Email: ", newRoot.Email)
	}

	router := middlewares.CorsMiddleware(r)

	log.Println("Server Started...")
	log.Fatal(http.ListenAndServe(port, router))
}
