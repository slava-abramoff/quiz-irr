package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"quiz-irr/internals/cache"
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

func getenvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func main() {
	// Optional .env for local development; in containers vars come from runtime env.
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using process environment")
	}

	secret := getenvOrDefault("SECRET", "FDGSSDFGSDFG")
	port := getenvOrDefault("BACKEND_PORT", ":8080")
	redisAddr := getenvOrDefault("REDIS_ADDR", "localhost:6379")
	root := dto.CreateUserRequest{
		FullName: getenvOrDefault("ADMIN_NAME", "Абрамов Вячеслав Александрович"),
		Password: getenvOrDefault("ADMIN_PASSWORD", "changeme"),
		Email:    getenvOrDefault("ADMIN_EMAIL", "vyachik005@gmail.com"),
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed connect: " + err.Error())
	}

	// Connect to redis; if unavailable, app still works without cache.
	redisClient, err := cache.ConnectCache(context.Background(), redisAddr)
	if err != nil {
		log.Println("Redis unavailable, caching disabled:", err)
		redisClient = nil
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
	examCases := usecases.NewExamCases(testsService, rawService, optionsService, questionsService, redisClient)
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
