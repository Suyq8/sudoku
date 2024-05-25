package Test

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"sudoku/pkg/sudoku"
	"testing"
)

func test(data []string, solverType string) error {
	board := make([]int32, 81)
	for i, c := range data[0] {
		board[i] = int32(c - '0')
	}

	solver := sudoku.SudokuSolver{Board: board}
	res, err := solver.Solve(solverType)

	if err != nil {
		return err
	}

	for i := 0; i < 81; i++ {
		if res[i] != int32(data[1][i]-'0') {
			log.Fatal(board, res, data[1])
			return errors.New("wrong answer for sudoku")
		}
	}
	return nil
}

func testSolver(t *testing.T, solverType string) {
	log.SetFlags(0)
	log.SetOutput(io.Discard)
	path := "data/sudoku.csv"
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	_, _ = reader.Read()

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}

		go func(data []string, solverType string) {
			board := make([]int32, 81)
			for i, c := range data[0] {
				board[i] = int32(c - '0')
			}

			solver := sudoku.SudokuSolver{Board: board}
			res, err := solver.Solve(solverType)

			if err != nil {
				t.Error(err)
			}

			for i := 0; i < 81; i++ {
				if res[i] != int32(data[1][i]-'0') {
					log.Fatal(board, res, data[1])
					t.Error("wrong answer for sudoku")
				}
			}
		}(data, solverType)
	}
}

func TestSolverBacktrack(t *testing.T) {
	testSolver(t, "backtrack")
}

func TestSolverDLX(t *testing.T) {
	testSolver(t, "dlx")
}
