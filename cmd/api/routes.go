package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	r := gin.New()

	r.Use(app.recoverPanic())
	r.Use(app.requestLogger())

	r.NoRoute(func(c *gin.Context) {
		app.notFoundResponse(c)
	})

	r.NoMethod(func(c *gin.Context) {
		app.methodNotAllowedResponse(c)
	})

	r.GET("/v1/healthcheck", app.healthCheckHandler)

	v1 := r.Group("/v1")
	{
		v1.POST("/students", app.createStudentHandler)
		v1.GET("/students", app.listStudentsHandler)
		v1.GET("/students/:id", app.showStudentHandler)
		v1.PATCH("/students/:id", app.updateStudentHandler)
		v1.DELETE("/students/:id", app.deleteStudentHandler)
	}

	return r

}
