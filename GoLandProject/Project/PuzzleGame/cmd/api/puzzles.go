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
	// Copy the values from the input struct to a new Movie struct.
	puzzle := &data.Puzzle{
		Title:        input.Title,
		NumOfPuzzles: input.NumOfPuzzles,
		Genres:       input.Genres,
	}
	// Initialize a new Validator.
	v := validator.New()
	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateMovie(v, puzzle); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showPuzzleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.
	puzzle := data.Puzzle{
		ID:           id,
		CreatedAt:    time.Now(),
		Title:        "Casablanca",
		NumOfPuzzles: 2000,
		Genres:       []string{"drama", "romance", "war"},
		Version:      1,
	}
	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, envelope{"puzzle": puzzle}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
