package main

import (
	"aoc/2023/common"
	"bufio"
)

type Direction byte

const (
	Up Direction = iota
	Right
	Down
	Left
)

type p struct {
	i, j int
}

type state struct {
	p p
	d Direction
}

const GRID_SIZE = 140

func main() {
	pipes, LRPipes, start := enhanceLowResPipes()

	// Find the cycle
	s := state{p: start, d: Up}
	s = next(pipes, s)

	count := 1
	for ; s.p != start; count++ {
		s = next(pipes, s)
	}

	// Mark enclosed areas
	for i, row := range pipes {
		for j, c := range row {
			if c == '.' {
				tiles, enclosed := floodFill(pipes, i, j, make([]p, 256))
				var filling = byte('O')
				if enclosed {
					filling = 'I'
				}
				for _, t := range tiles {
					pipes[t.i][t.j] = filling
				}
			}
		}
	}

	// Count enclosed areas
	var nestSpots int
	for i, row := range LRPipes {
		for j := range row {
			if pipes[2*i][2*j] == 'I' {
				nestSpots++
			}
		}
	}
	println(count / 4)
	println(nestSpots)
}

func floodFill(pipes [][]byte, i, j int, tiles []p) ([]p, bool) {
	// Base case: visited / blocked tile / might be enclosed
	curr := pipes[i][j]
	if curr == '#' {
		return tiles, true
	}

	// Add visited
	pipes[i][j] = '#'
	tiles = append(tiles, p{i: i, j: j})

	// Base case: is not enclosed
	if i == 0 || i == (2*GRID_SIZE)-1 || j == 0 || j == (2*GRID_SIZE)-1 {
		return tiles, false
	}

	// Recurse
	var enclosed, temp bool
	tiles, enclosed = floodFill(pipes, i+1, j, tiles)
	tiles, temp = floodFill(pipes, i-1, j, tiles)
	enclosed = enclosed && temp
	tiles, temp = floodFill(pipes, i, j+1, tiles)
	enclosed = enclosed && temp
	tiles, temp = floodFill(pipes, i, j-1, tiles)
	return tiles, enclosed && temp
}

func enhanceLowResPipes() ([][]byte, []string, p) {
	file := common.Open("input")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	LRPipes := make([]string, 0, 140)
	pipes := make([][]byte, 2*GRID_SIZE)
	for i := range pipes {
		pipes[i] = make([]byte, 2*GRID_SIZE)
	}

	var start p
	for i := 0; scanner.Scan(); i++ {
		row := scanner.Text()
		LRPipes = append(LRPipes, row)
		s := zoomIn(pipes, row, i)
		if s.i != 0 {
			start = s
		}
	}

	return pipes, LRPipes, start
}

func zoomIn(pipes [][]byte, row string, i int) p {
	var start p
	for j, c := range row {
		pipes[2*i+1][2*j+1] = '.'
		switch c {
		case '|':
			pipes[2*i][2*j] = '|'
			pipes[2*i][2*j+1] = '.'
			pipes[2*i+1][2*j] = '|'
		case '-':
			pipes[2*i][2*j] = '-'
			pipes[2*i][2*j+1] = '-'
			pipes[2*i+1][2*j] = '.'
		case 'L':
			pipes[2*i][2*j] = 'L'
			pipes[2*i][2*j+1] = '-'
			pipes[2*i+1][2*j] = '.'
		case 'J':
			pipes[2*i][2*j] = 'J'
			pipes[2*i][2*j+1] = '.'
			pipes[2*i+1][2*j] = '.'
		case '7':
			pipes[2*i][2*j] = '7'
			pipes[2*i][2*j+1] = '.'
			pipes[2*i+1][2*j] = '|'
		case 'F':
			pipes[2*i][2*j] = 'F'
			pipes[2*i][2*j+1] = '-'
			pipes[2*i+1][2*j] = '|'
		case '.':
			pipes[2*i][2*j] = '.'
			pipes[2*i][2*j+1] = '.'
			pipes[2*i+1][2*j] = '.'
		case 'S':
			start = p{i: 2 * i, j: 2 * j}
			pipes[2*i][2*j] = '|'
			pipes[2*i][2*j+1] = '.'
			pipes[2*i+1][2*j] = '|'
		}
	}
	return start
}

func next(pipes [][]byte, s state) state {
	temp := s.p
	switch pipes[s.p.i][s.p.j] {
	case '|':
		if s.d == Up {
			s.p.i--
		} else {
			s.p.i++
		}
	case '-':
		if s.d == Right {
			s.p.j++
		} else {
			s.p.j--
		}
	case 'L':
		if s.d == Down {
			s.d = Right
			s.p.j++
		} else {
			s.d = Up
			s.p.i--
		}
	case 'J':
		if s.d == Down {
			s.d = Left
			s.p.j--
		} else {
			s.d = Up
			s.p.i--
		}
	case '7':
		if s.d == Right {
			s.d = Down
			s.p.i++
		} else {
			s.d = Left
			s.p.j--
		}
	case 'F':
		if s.d == Left {
			s.d = Down
			s.p.i++
		} else {
			s.d = Right
			s.p.j++
		}
	}
	pipes[temp.i][temp.j] = '#'
	return s
}
