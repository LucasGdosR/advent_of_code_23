package main

import (
	"aoc/2023/common"
	"bufio"
)

type state struct {
	i, j, dir int
}

const (
	LENGTH    = 110
	SENTINEL  = 1
	SENTINELS = 2 * SENTINEL
	UP        = iota
	RIGHT
	DOWN
	LEFT
)

var (
	grid      [][]byte
	seen      = map[state]struct{}{}
	energized = make([][]bool, LENGTH+SENTINELS)
)

func main() {
	initGrid()
	start := state{i: SENTINEL, j: 0, dir: RIGHT}

	for i := range energized {
		energized[i] = make([]bool, LENGTH+SENTINELS)
	}

	propagateBeam(start)

	var count int
	for i := SENTINEL; i < LENGTH+SENTINEL; i++ {
		for j := SENTINEL; j < LENGTH+SENTINEL; j++ {
			if energized[i][j] {
				count++
			}
		}
	}
	println(count)
}

func propagateBeam(s state) {
	next := getNextTile(s)
	if _, ok := seen[next]; ok {
		return
	}

	seen[next] = struct{}{}
	energized[next.i][next.j] = true

	switch grid[next.i][next.j] {
	case '.':
		propagateBeam(next)
	case '\\':
		switch next.dir {
		case UP:
			next.dir = LEFT
		case RIGHT:
			next.dir = UP
		case DOWN:
			next.dir = RIGHT
		case LEFT:
			next.dir = DOWN
		}
		propagateBeam(next)
	case '/':
		switch next.dir {
		case UP:
			next.dir = RIGHT
		case RIGHT:
			next.dir = UP
		case DOWN:
			next.dir = LEFT
		case LEFT:
			next.dir = DOWN
		}
		propagateBeam(next)
	case '-':
		if next.dir == LEFT || next.dir == RIGHT {
			propagateBeam(next)
		} else {
			next.dir = LEFT
			propagateBeam(next)
			next.dir = RIGHT
			propagateBeam(next)
		}
	case '|':
		if next.dir == UP || next.dir == DOWN {
			propagateBeam(next)
		} else {
			next.dir = UP
			propagateBeam(next)
			next.dir = DOWN
			propagateBeam(next)
		}
	}
}

func getNextTile(s state) state {
	switch s.dir {
	case UP:
		s.i--
	case RIGHT:
		s.j++
	case DOWN:
		s.i++
	case LEFT:
		s.j--
	}
	return s
}

func initGrid() [][]byte {
	grid = make([][]byte, LENGTH+SENTINELS)
	for i := range grid {
		grid[i] = make([]byte, LENGTH+SENTINELS)
	}

	file := common.Open("input")
	scanner := bufio.NewScanner(file)

	for i := 1; scanner.Scan(); i++ {
		line := scanner.Text()
		copy(grid[i][1:], line)
	}

	return grid
}
