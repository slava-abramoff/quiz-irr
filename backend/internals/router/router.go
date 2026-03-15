package router

import (
	"quiz-irr/internals/handlers"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(
	a *handlers.AuthHandlers,
	u *handlers.UserHandlers,
	t *handlers.TestHandlers,
	e *handlers.ExamHandlers,
	auth func(onlyRoot bool) func(httprouter.Handle) httprouter.Handle,
) *httprouter.Router {
	router := httprouter.New()

	basicAuth := auth(false)
	onlyRoot := auth(true)
	// users
	router.POST("/api/users", onlyRoot(u.Create))
	router.PATCH("/api/users/:id", onlyRoot(u.Update))
	router.DELETE("/api/users/:id", onlyRoot(u.Delete))
	router.POST("/api/auth/login", a.Login)
	router.POST("/api/auth/refresh", a.Refresh)

	// Работа админа с тестами
	router.POST("/api/tests/init", basicAuth(t.NewTest))
	router.GET("/api/tests/:id/:mode", t.GetTest)
	router.GET("/api/tests", t.FindManyTests)
	router.PATCH("/api/tests/:id", t.UpdateTest)
	router.DELETE("/api/tests/:id", t.DeleteTest)

	router.POST("/api/questions/test/:id", t.AddQuestion)
	router.PATCH("/api/questions/:id", t.EditQuestion)
	router.DELETE("/api/questions/:id", t.DeleteQuestion)

	router.POST("/api/options/question/:id", t.AddOption)
	router.PATCH("/api/options/:id", t.EditOption)
	router.DELETE("/api/options/:id", t.DeleteOption)

	// Прохождение теста
	router.GET("/api/exam/info/:testId", e.GetTestInfo)
	router.POST("/api/exam/start/:testId", e.StartTest)
	router.POST("/api/exam/save/:rawId", e.SaveAnswers)

	return router
}
