package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/puzzles", app.requirePermission("puzzles:read", app.listPuzzlesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/puzzles", app.requirePermission("puzzles:write", app.createPuzzleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/puzzles/:id", app.requirePermission("puzzles:read", app.showPuzzleHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/puzzles/:id", app.requirePermission("puzzles:write", app.updatePuzzleHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/puzzles/:id", app.requirePermission("puzzles:write", app.deletePuzzleHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
