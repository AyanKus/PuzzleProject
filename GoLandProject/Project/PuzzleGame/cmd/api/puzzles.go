package main

import (
	"Puzzle.Ayan.net/internal/data"
	"Puzzle.Ayan.net/internal/validator"
	"errors"
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

	err = app.models.Puzzles.Insert(puzzle)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/puzzles/%d", puzzle.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"puzzle": puzzle}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPuzzleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	puzzle, err := app.models.Puzzles.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"puzzle": puzzle}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie := data.Puzzle{
		ID:           id,
		CreatedAt:    time.Now(),
		Title:        "Casablanca",
		NumOfPuzzles: 102,
		Genres:       []string{"drama", "romance", "war"},
		Version:      1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updatePuzzleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	puzzle, err := app.models.Puzzles.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title        *string   `json:"title"`
		NumOfPuzzles *data.NOP `json:"num_of_puzzles"`
		Genres       []string  `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Title != nil {
		puzzle.Title = *input.Title
	}
	if input.NumOfPuzzles != nil {
		puzzle.NumOfPuzzles = *input.NumOfPuzzles
	}
	if input.Genres != nil {
		puzzle.Genres = input.Genres
	}
	v := validator.New()
	if data.ValidateMovie(v, puzzle); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Puzzles.Update(puzzle)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"puzzle": puzzle}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePuzzleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Puzzles.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) listPuzzlesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "num_of_puzzles", "-id", "-title", "-num_of_puzzles"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	puzzles, metadata, err := app.models.Puzzles.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Include the metadata in the response envelope.
	err = app.writeJSON(w, http.StatusOK, envelope{"puzzles": puzzles, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
