package main

import (
	"fmt"
	"net/http"
)

// Log and handle generic errors
func (a *applicationDependencies) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	a.logger.Error(err.Error(), "method", method, "uri", uri)
}

// Send a JSON error response
func (a *applicationDependencies) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message any) {
	errorData := envelope{"error": message}
	err := a.writeJSON(w, status, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Send a 500 Internal Server Error response
func (a *applicationDependencies) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	a.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

// Send a 404 Not Found response with a custom message
func (a *applicationDependencies) notFoundResponse(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "the requested resource could not be found"
	}
	a.errorResponseJSON(w, r, http.StatusNotFound, message)
}

// Send a 405 Method Not Allowed response
func (a *applicationDependencies) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponseJSON(w, r, http.StatusMethodNotAllowed, message)
}

// Send a 400 Bad Request response with a custom error message
func (a *applicationDependencies) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}

// Send a 422 Unprocessable Entity response for validation errors
func (a *applicationDependencies) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponseJSON(w, r, http.StatusUnprocessableEntity, errors)
}

func (a *applicationDependencies) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	a.errorResponseJSON(w, r, http.StatusTooManyRequests, message)
}
