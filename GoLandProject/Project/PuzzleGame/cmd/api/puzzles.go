package main

import (
	"Puzzle.Ayan.net/internal/data"
	"Puzzle.Ayan.net/internal/validator"
	"errors"
	"fmt"
	"net/http"
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
	game, err := app.models.Puzzles.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
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
		Title        string   `json:"title"`
		NumOfPuzzles data.NOP `json:"num_of_puzzles"`
		Genres       []string `json:"genres"`
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

func (app *application) deleteGameHandler(w http.ResponseWriter, r *http.Request) {
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
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "game successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listGamesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Games []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.Games = app.readCSV(qs, "games", []string{})
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "score", "-id", "-title", "-score"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	games, metadata, err := app.models.Puzzles.GetAll(input.Title, input.Games, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"games": games, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
