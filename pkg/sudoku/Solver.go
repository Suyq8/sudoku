package sudoku

import (
	context "context"
	"errors"
	"log"
	"time"
)

type Solver interface {
	Solve() ([]int32, error)
}

type SudokuSolver struct {
	Board []int32
	Solver
	UnimplementedSudokuServer
}

func (s *SudokuSolver) SolverFactory(solverType string) (Solver, error) {
	switch solverType {
	case "backtrack":
		return &SolverBacktrack{Board: s.Board}, nil
	case "dlx":
		return &SolverDLX{Board: s.Board}, nil
	default:
		return nil, errors.New("invalid solver type")
	}
}

func (s *SudokuSolver) CheckNumberValid(row int, col int, num int32) bool {
	// Check row, column
	for i := 0; i < 9; i++ {
		if s.Board[GetIndex(row, i)] == num || s.Board[GetIndex(i, col)] == num {
			return false
		}
	}

	// Check 3x3 box
	startRow := row / 3 * 3
	startCol := col / 3 * 3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.Board[GetIndex(startRow+i, startCol+j)] == num {
				return false
			}
		}
	}

	return true
}

func (s *SudokuSolver) Solve(solverType string) ([]int32, error) {
	// check whether the board is valid
	if len(s.Board) != 81 {
		return []int32{}, errors.New("size of board is not 81")
	}

	/*for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.Board[GetIndex(i, j)] < 0 || s.Board[GetIndex(i, j)] > 9 {
				return []int32{}, errors.New("invalid number in board")
			}
			if s.Board[GetIndex(i, j)] != 0 {
				tmp := s.Board[GetIndex(i, j)]
				s.Board[GetIndex(i, j)] = 0
				if !s.CheckNumberValid(i, j, tmp) {
					return []int32{}, errors.New("invalid board")
				}
				s.Board[GetIndex(i, j)] = tmp
			}
		}
	}*/

	solver, err := s.SolverFactory(solverType)
	if err != nil {
		return []int32{}, err
	}

	return solver.Solve()
}

func (s *SudokuSolver) GetSolution(ctx context.Context, in *Question) (*Solution, error) {
	board := in.GetBoard()
	s.Board = board
	t := time.Now()
	res, err := s.Solve(in.GetSolverType())
	duration := time.Since(t)
	if err != nil {
		log.Println(err)
	}

	return &Solution{Board: res, Solved: err == nil, Duration: duration.Seconds()}, nil
}

func NewSolver() *SudokuSolver {
	return &SudokuSolver{}
}
