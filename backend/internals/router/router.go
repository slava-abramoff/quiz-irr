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
	rw *handlers.RawAnswersHandlers,
	res *handlers.ResultsHandlers,
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
	router.GET("/api/tests/:id/:mode", basicAuth(t.GetTest))
	router.GET("/api/tests", basicAuth(t.FindManyTests))
	router.PATCH("/api/tests/:id", basicAuth(t.UpdateTest))
	router.DELETE("/api/tests/:id", basicAuth(t.DeleteTest))

	router.POST("/api/questions/test/:id", basicAuth(t.AddQuestion))
	router.PATCH("/api/questions/:id", basicAuth(t.EditQuestion))
	router.DELETE("/api/questions/:id", basicAuth(t.DeleteQuestion))

	router.POST("/api/options/question/:id", basicAuth(t.AddOption))
	router.PATCH("/api/options/:id", basicAuth(t.EditOption))
	router.DELETE("/api/options/:id", basicAuth(t.DeleteOption))

	// Прохождение теста
	router.GET("/api/exam/info/:testId", e.GetTestInfo)
	router.POST("/api/exam/start/:testId", e.StartTest)
	router.POST("/api/exam/save/:rawId", e.SaveAnswers)

	// Сырые результаты
	// router.GET("/api/raws/test/:testId", basicAuth(rw.FindRawResults))
	router.GET("/api/raws/test/:testId", rw.FindRawResults)
	router.POST("/api/raws/analyze/:rawId", basicAuth(rw.AnalyzeResults))
	// POST /api/raws/test/analyze/:testId
	// GET /api/raws/answers/report/:rawId
	router.DELETE("/api/raws/answers/:rawId", basicAuth(rw.DeleteRawResults))
	router.DELETE("/api/raws/test/:testId", basicAuth(rw.DeleteAllRawByTest))

	// Результаты и рейтинг
	router.GET("/api/results/test/:testId", basicAuth(res.GetListByTest))
	router.DELETE("/api/results/result/:resultId", basicAuth(res.DeleteResult))
	router.DELETE("/api/results/test/:testId", basicAuth(res.DeleteResultsByTest))

	return router
}
