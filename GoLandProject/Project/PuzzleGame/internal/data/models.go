package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Puzzles PuzzleModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Puzzles: PuzzleModel{DB: db},
	}
}
