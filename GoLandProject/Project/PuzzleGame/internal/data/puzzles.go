package data

import (
	"time"
)

type Puzzle struct {
	ID           int64        `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	Title        string       `json:"title"`
	NumOfPuzzles numOfPuzzles `json:"num_of_puzzles,omitempty,string"`
	Genres       []string     `json:"genres,omitempty"`
	Version      int32        `json:"version"`
}
