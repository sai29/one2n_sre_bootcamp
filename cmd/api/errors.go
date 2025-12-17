package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// logError logs an error with request context information.
func (app *application) logError(c *gin.Context, err error) {
	app.logger.PrintError(err, map[string]string{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	})
}

// errorResponse sends a JSON error response.
func (app *application) errorResponse(c *gin.Context, status int, message interface{}) {
	c.JSON(status, gin.H{"error": message})
}

// === 5xx Errors ===

func (app *application) serverErrorResponse(c *gin.Context, err error) {
	app.logError(c, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(c, http.StatusInternalServerError, message)
}

// === 4xx Errors ===

func (app *application) notFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	app.errorResponse(c, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	app.errorResponse(c, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}

//nolint:unused
func (app *application) failedValidationResponse(c *gin.Context, errors map[string]string) {
	app.errorResponse(c, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(c *gin.Context) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(c, http.StatusConflict, message)
}

//nolint:unused
func (app *application) rateLimitExceededResponse(c *gin.Context) {
	message := "rate limit exceeded"
	app.errorResponse(c, http.StatusTooManyRequests, message)
}

// === Authentication / Authorization Errors ===
//
//nolint:unused
func (app *application) invalidCredentialsResponse(c *gin.Context) {
	message := "invalid authentication credentials"
	app.errorResponse(c, http.StatusUnauthorized, message)
}

//nolint:unused
func (app *application) invalidAuthenticationTokenResponse(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(c, http.StatusUnauthorized, message)
}

//nolint:unused
func (app *application) authenticationRequiredResponse(c *gin.Context) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(c, http.StatusUnauthorized, message)
}

//nolint:unused
func (app *application) inactiveAccountResponse(c *gin.Context) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(c, http.StatusForbidden, message)
}

//nolint:unused
func (app *application) notPermittedResponse(c *gin.Context) {
	message := "you do not have the necessary permissions to access this resource"
	app.errorResponse(c, http.StatusForbidden, message)
}
