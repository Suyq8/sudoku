package sudoku

import (
	"errors"
	"log"
)

type SolverBacktrack struct {
	Board    []int32
	MaskCol  [9]uint
	MaskRow  [9]uint
	MaskGrid [9]uint
}

func GetIndex(row int, col int) int {
	return row*9 + col
}

func (s *SolverBacktrack) CheckNumberValid(row int, col int, num int32) bool {
	return s.MaskCol[col]&(1<<num) == 0 && s.MaskRow[row]&(1<<num) == 0 && s.MaskGrid[(row/3)*3+col/3]&(1<<num) == 0
}

func (s *SolverBacktrack) SetNumber(row int, col int, num int32) {
	s.Board[GetIndex(row, col)] = num
	s.MaskRow[row] |= 1 << num
	s.MaskCol[col] |= 1 << num
	s.MaskGrid[(row/3)*3+col/3] |= 1 << num
}

func (s *SolverBacktrack) RemoveNumber(row int, col int, num int32) {
	s.Board[GetIndex(row, col)] = 0
	s.MaskRow[row] &= ^(1 << num)
	s.MaskCol[col] &= ^(1 << num)
	s.MaskGrid[(row/3)*3+col/3] &= ^(1 << num)
}

func (s *SolverBacktrack) SolveBacktrack() bool {
	nextRow, nextCol, hasEmptyCell := s.FindEmptyCell()
	if !hasEmptyCell {
		return true
	}

	for candidate := int32(1); candidate <= 9; candidate++ {
		if s.CheckNumberValid(nextRow, nextCol, candidate) {
			s.SetNumber(nextRow, nextCol, candidate)

			if s.SolveBacktrack() {
				return true
			}
			s.RemoveNumber(nextRow, nextCol, candidate)
		}
	}

	return false
}

func (s *SolverBacktrack) FindEmptyCell() (int, int, bool) {
	for i := 0; i < 9; i++ {
		if s.MaskRow[i] == 0x3ff {
			continue
		}
		for j := 0; j < 9; j++ {
			if s.Board[GetIndex(i, j)] == 0 {
				return i, j, true
			}
		}
	}

	return 0, 0, false
}

func (s *SolverBacktrack) Solve() ([]int32, error) {
	log.Println("Solving Sudoku using Backtrack")
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			num := s.Board[GetIndex(i, j)]
			if num < 0 || num > 9 {
				return []int32{}, errors.New("invalid number in board")
			}

			if num != 0 {
				if s.MaskRow[i]&(1<<num) != 0 {
					return []int32{}, errors.New("invalid number in row")
				} else {
					s.MaskRow[i] |= 1 << num
				}

				if s.MaskCol[j]&(1<<num) != 0 {
					return []int32{}, errors.New("invalid number in col")
				} else {
					s.MaskCol[j] |= 1 << num
				}

				if s.MaskGrid[(i/3)*3+j/3]&(1<<num) != 0 {
					return []int32{}, errors.New("invalid number in grid")
				} else {
					s.MaskGrid[(i/3)*3+j/3] |= 1 << num
				}
			}
		}
	}

	if !s.SolveBacktrack() {
		return []int32{}, errors.New("failed")
	}
	log.Println("Solved")

	return s.Board, nil
}
