package main

import (
	"Puzzle.Ayan.net/internal/data"
	"Puzzle.Ayan.net/internal/validator"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createPuzzleHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title        string   `json:"title"`
		NumOfPuzzles data.NOP `json:"num_of_puzzles"`
		Genres       []string `json:"genres"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	puzzle := &data.Puzzle{
		Title:        input.Title,
		NumOfPuzzles: input.NumOfPuzzles,
		Genres:       input.Genres,
	}
	v := validator.New()
	if data.ValidateMovie(v, puzzle); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showPuzzleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	puzzle := data.Puzzle{
		ID:           id,
		CreatedAt:    time.Now(),
		Title:        "Casablanca",
		NumOfPuzzles: 2000,
		Genres:       []string{"drama", "romance", "war"},
		Version:      1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"puzzle": puzzle}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
