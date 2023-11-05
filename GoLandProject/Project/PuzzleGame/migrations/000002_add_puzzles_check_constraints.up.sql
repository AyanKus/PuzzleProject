ALTER TABLE puzzles ADD CONSTRAINT puzzles_NumOfPuzzles_check CHECK ( NOP >= 0);
ALTER TABLE puzzles ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
