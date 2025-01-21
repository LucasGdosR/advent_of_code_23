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

func main() {
	file := common.Open("input")
	scanner := bufio.NewScanner(file)
	pipes := make([][]byte, 0, 140)
	for scanner.Scan() {
		bytes := make([]byte, 140)
		copy(bytes, scanner.Bytes())
		pipes = append(pipes, bytes)
	}

	start := findStart(pipes)

	s := state{p: start, d: Up}
	s = next(pipes, s)

	count := 1
	for ; s.p != start; count++ {
		s = next(pipes, s)
	}

	var nestSpots int
	for _, row := range pipes {
		var isEnclosed bool
		for _, c := range row {
			if c == '#' {
				isEnclosed = !isEnclosed
			} else if isEnclosed && c != '@' {
				nestSpots++
			}
		}
	}

	println(count / 2)

	println(nestSpots)
}

func next(pipes [][]byte, s state) state {
	temp := s.p
	c := byte('#')
	switch pipes[s.p.i][s.p.j] {
	case '|':
		if s.d == Up {
			s.p.i--
		} else {
			s.p.i++
		}
	case '-':
		c = '@'
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
		c = '@'
		if s.d == Right {
			s.d = Down
			s.p.i++
		} else {
			s.d = Left
			s.p.j--
		}
	case 'F':
		c = '@'
		if s.d == Left {
			s.d = Down
			s.p.i++
		} else {
			s.d = Right
			s.p.j++
		}
	case 'S':
		pipes[s.p.i][s.p.j] = '#'
		s.p.i--
	}
	pipes[temp.i][temp.j] = c
	return s
}

func findStart(pipes [][]byte) p {
	for i, row := range pipes {
		for j, c := range row {
			if c == 'S' {
				return p{i: i, j: j}
			}
		}
	}
	return p{i: -1, j: -1}
}
