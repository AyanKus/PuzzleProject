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
	router.HandlerFunc(http.MethodGet, "/v1/puzzles", app.listPuzzlesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/puzzles", app.createPuzzleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/puzzles/:id", app.showPuzzleHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/puzzles/:id", app.updatePuzzleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/puzzles/:id", app.deletePuzzleHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	return app.recoverPanic(app.rateLimit(router))
}
