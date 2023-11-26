package data

import (
	"Puzzle.Ayan.net/internal/validator"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
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

type PuzzleModel struct {
	DB *sql.DB
}

func (m PuzzleModel) Insert(puzzle *Puzzle) error {
	query := `
		INSERT INTO puzzles (title, NumOfPuzzles, genres)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`
	args := []interface{}{puzzle.Title, puzzle.NumOfPuzzles, pq.Array(puzzle.Genres)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&puzzle.ID, &puzzle.CreatedAt, &puzzle.Version)
}

func (m PuzzleModel) Get(id int64) (*Puzzle, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, created_at, numOfPuzzles, genres, version
		FROM puzzles
		WHERE id = $1`
	var puzzle Puzzle
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&puzzle.ID,
		&puzzle.CreatedAt,
		&puzzle.Title,
		&puzzle.NumOfPuzzles,
		pq.Array(&puzzle.Genres),
		&puzzle.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &puzzle, nil
}

func (m PuzzleModel) Update(puzzle *Puzzle) error {
	query := `
		UPDATE puzzles
		SET title = $1, numOfPuzzles = $2, genres = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`
	args := []interface{}{
		puzzle.Title,
		puzzle.NumOfPuzzles,
		pq.Array(puzzle.Genres),
		puzzle.ID,
		puzzle.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&puzzle.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m PuzzleModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM puzzles
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m PuzzleModel) GetAll(title string, genres []string, filters Filters) ([]*Puzzle, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT id, created_at, title, numOfPuzzles, genres, version
		FROM puzzles
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()
	totalRecords := 0
	puzzles := []*Puzzle{}
	for rows.Next() {
		var puzzle Puzzle
		err := rows.Scan(
			&totalRecords,
			&puzzle.ID,
			&puzzle.CreatedAt,
			&puzzle.Title,
			&puzzle.NumOfPuzzles,
			pq.Array(&puzzle.Genres),
			&puzzle.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		puzzles = append(puzzles, &puzzle)
	}
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err // Update this to return an empty Metadata struct.
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return puzzles, metadata, nil
}
