package router

import (
	"quiz-irr/internals/handlers"

	"github.com/julienschmidt/httprouter"
)

func NewRouter(
	u *handlers.UserHandlers,
) *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/users", u.Create)
	router.PATCH("/api/users/:id", u.Update)
	router.DELETE("/api/users/:id", u.Delete)

	return router
}
