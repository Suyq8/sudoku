package sudoku

import (
	"errors"
	"log"
)

type SolverDLX struct {
	Board    []int32
	solver   DLX
	MaskCol  [9]uint
	MaskRow  [9]uint
	MaskGrid [9]uint
}

const N = 5e4
const MAXSIZE = 1e4

type DLX struct {
	n     int
	m     int
	tot   int
	first [MAXSIZE + 10]int
	siz   [MAXSIZE + 10]int
	L     [MAXSIZE + 10]int
	R     [MAXSIZE + 10]int
	U     [MAXSIZE + 10]int
	D     [MAXSIZE + 10]int
	col   [MAXSIZE + 10]int
	row   [MAXSIZE + 10]int
	ans   [9][9]int
	stk   [N]int
}

func (d *DLX) build(r, c int) {
	d.n = r
	d.m = c
	for i := 0; i <= c; i++ {
		d.L[i] = i - 1
		d.R[i] = i + 1
		d.U[i] = i
		d.D[i] = i
	}
	d.L[0] = c
	d.R[c] = 0
	d.tot = c
}

func (d *DLX) insert(r, c int) {
	d.tot++
	d.col[d.tot] = c
	d.row[d.tot] = r
	d.siz[c]++
	d.D[d.tot] = d.D[c]
	d.U[d.D[c]] = d.tot
	d.U[d.tot] = c
	d.D[c] = d.tot
	if d.first[r] == 0 {
		d.first[r] = d.tot
		d.L[d.tot] = d.tot
		d.R[d.tot] = d.tot
	} else {
		d.R[d.tot] = d.R[d.first[r]]
		d.L[d.R[d.first[r]]] = d.tot
		d.L[d.tot] = d.first[r]
		d.R[d.first[r]] = d.tot
	}
}

func (d *DLX) remove(c int) {
	var i, j int
	d.L[d.R[c]] = d.L[c]
	d.R[d.L[c]] = d.R[c]
	for i = d.D[c]; i != c; i = d.D[i] {
		for j = d.R[i]; j != i; j = d.R[j] {
			d.U[d.D[j]] = d.U[j]
			d.D[d.U[j]] = d.D[j]
			d.siz[d.col[j]]--
		}
	}
}

func (d *DLX) recover(c int) {
	var i, j int
	for i = d.U[c]; i != c; i = d.U[i] {
		for j = d.L[i]; j != i; j = d.L[j] {
			d.U[d.D[j]] = j
			d.D[d.U[j]] = j
			d.siz[d.col[j]]++
		}
	}
	d.L[d.R[c]] = c
	d.R[d.L[c]] = c
}

func (d *DLX) dance(dep int) bool {
	c := d.R[0]
	if d.R[0] == 0 {
		for i := 1; i < dep; i++ {
			x := (d.stk[i] - 1) / 9 / 9
			y := (d.stk[i] - 1) / 9 % 9
			v := (d.stk[i]-1)%9 + 1
			d.ans[x][y] = v
		}
		return true
	}
	for i := d.R[0]; i != 0; i = d.R[i] {
		if d.siz[i] < d.siz[c] {
			c = i
		}
	}
	d.remove(c)
	for i := d.D[c]; i != c; i = d.D[i] {
		d.stk[dep] = d.row[i]
		for j := d.R[i]; j != i; j = d.R[j] {
			d.remove(d.col[j])
		}
		if d.dance(dep + 1) {
			return true
		}
		for j := d.L[i]; j != i; j = d.L[j] {
			d.recover(d.col[j])
		}
	}
	d.recover(c)
	return false
}

func (s *SolverDLX) GetId(row, col, num int) int {
	return (row-1)*9*9 + (col-1)*9 + num
}

func (s *SolverDLX) Insert(row, col, num int) {
	dx := (row-1)/3 + 1
	dy := (col-1)/3 + 1
	room := (dx-1)*3 + dy
	id := s.GetId(row, col, num)
	f1 := (row-1)*9 + num         // task 1
	f2 := 81 + (col-1)*9 + num    // task 2
	f3 := 81*2 + (room-1)*9 + num // task 3
	f4 := 81*3 + (row-1)*9 + col  // task 4
	s.solver.insert(id, f1)
	s.solver.insert(id, f2)
	s.solver.insert(id, f3)
	s.solver.insert(id, f4)
}

func (s *SolverDLX) Solve() ([]int32, error) {
	log.Println("Solving Sudoku using Dancing Links X")
	s.solver.build(729, 324)

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			s.solver.ans[i][j] = int(s.Board[i*9+j])
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

			for v := 1; v <= 9; v++ {
				if s.solver.ans[i][j] != 0 && s.solver.ans[i][j] != v {
					continue
				}
				s.Insert(i+1, j+1, v)
			}
		}
	}

	if solved := s.solver.dance(1); !solved {
		return []int32{}, errors.New("failed")
	}
	log.Println("Solved")

	var res []int32
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			res = append(res, int32(s.solver.ans[i][j]))
		}
	}
	return res, nil
}
