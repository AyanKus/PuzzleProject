package data

import (
	"Puzzle.Ayan.net/internal/validator"
	"time"
)

type Puzzle struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Title        string    `json:"title"`
	NumOfPuzzles NOP       `json:"num_of_puzzles,omitempty,string"`
	Genres       []string  `json:"genres,omitempty"`
	Version      int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, puzzle *Puzzle) {
	v.Check(puzzle.Title != "", "title", "must be provided")
	v.Check(len(puzzle.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(puzzle.NumOfPuzzles != 0, "number of puzzles", "must be provided")
	v.Check(puzzle.NumOfPuzzles > 0, "number of puzzles", "must be a positive integer")
	v.Check(puzzle.Genres != nil, "genres", "must be provided")
	v.Check(len(puzzle.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(puzzle.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(puzzle.Genres), "genres", "must not contain duplicate values")
}
