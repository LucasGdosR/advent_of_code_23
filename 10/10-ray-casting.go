package main

import (
	"aoc/2023/common"
	"bufio"
	"runtime"
)

type Direction byte

const (
	GRID_SIDE           = 140
	Up        Direction = iota
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
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  findCycleAndNestSpotsST,
		MT_Impl:  findCycleAndNestSpotsMT,
		Part1Str: "Steps to the middle of the cycle",
		Part2Str: "Possible nest spots",
	}
	common.Benchmark(thisProgram, 1000)
}

func findCycleAndNestSpotsST() common.Results[int, int] {
	pipes := readPipeMaze()
	return common.Results[int, int]{
		Part1: findCycle(pipes) / 2,
		Part2: findNestSpots(pipes)}
}

func findCycleAndNestSpotsMT() common.Results[int, int] {
	pipes := readPipeMaze()
	cycleLength := findCycle(pipes)

	numWorkers := runtime.GOMAXPROCS(0)
	linesPerWoker := (GRID_SIDE + numWorkers - 1) / numWorkers
	partialNestSpots := make(chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * linesPerWoker
		end := start + linesPerWoker
		if i == numWorkers-1 {
			end = GRID_SIDE
		}
		go func(start, end int) {
			partialNestSpots <- findNestSpots(pipes[start:end])
		}(start, end)
	}

	results := common.Results[int, int]{Part1: cycleLength / 2}
	for i := 0; i < numWorkers; i++ {
		results.Part2 += <-partialNestSpots
	}
	close(partialNestSpots)

	return results
}

func findCycle(pipes [][]byte) int {
	start := findStart(pipes)

	s := state{p: start, d: Up}
	s = next(pipes, s)

	count := 1
	for ; s.p != start; count++ {
		s = next(pipes, s)
	}
	return count
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

func findNestSpots(pipes [][]byte) int {
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
	return nestSpots
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

func readPipeMaze() [][]byte {
	file := common.Open("input")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pipes := make([][]byte, 0, GRID_SIDE)
	for scanner.Scan() {
		bytes := make([]byte, GRID_SIDE)
		copy(bytes, scanner.Bytes())
		pipes = append(pipes, bytes)
	}
	return pipes
}
