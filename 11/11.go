package main

import (
	"aoc/2023/common"
	"syscall"
)

type p struct {
	i, j int
}

const (
	UNIVERSE_LENGTH = 140
	LINE_BREAK      = 1
)

func main() {
	var galaxiesPt1, galaxiesPt2 []p
	galaxiesPt1, galaxiesPt2 = expandUniverse(readUniverse())

	var results common.Results[int, int]
	for i, thisG := range galaxiesPt1 {
		for j, thatG := range galaxiesPt1[i+1:] {
			results.Part1 += abs(thisG, thatG)
			results.Part2 += abs(galaxiesPt2[i], galaxiesPt2[i+j+1])
		}
	}
	println(results.Part1, results.Part2)
}

func readUniverse() ([]int, []int, []p) {
	mappedFile := common.Mmap("input")
	universe := mappedFile.File
	defer syscall.Munmap(universe)

	emptyRows := make([]int, UNIVERSE_LENGTH)
	emptyCols := make([]int, UNIVERSE_LENGTH)
	for i := range emptyCols {
		emptyCols[i] = 1
	}
	galaxies := make([]p, 0, 443)

	for emptyRowCount, i := 0, 0; i < UNIVERSE_LENGTH; i++ {
		emptyRow := true
		for j := 0; j < UNIVERSE_LENGTH; j++ {
			if universe[(UNIVERSE_LENGTH+LINE_BREAK)*i+j] == '#' {
				galaxies = append(galaxies, p{i: i, j: j})
				emptyRow = false
				emptyCols[j] = 0
			}
		}
		if emptyRow {
			emptyRowCount++
		}
		emptyRows[i] = emptyRowCount
	}

	emptyColCount := 0
	for i, empty := range emptyCols {
		emptyColCount += empty
		emptyCols[i] = emptyColCount
	}

	return emptyRows, emptyCols, galaxies
}

func expandUniverse(emptyRows, emptyCols []int, galaxies []p) ([]p, []p) {
	galaxiesPt1 := make([]p, len(galaxies))
	galaxiesPt2 := make([]p, len(galaxies))

	for i, g := range galaxies {
		er, ec := emptyRows[g.i], emptyCols[g.j]
		galaxiesPt1[i] = p{i: g.i + er, j: g.j + ec}
		galaxiesPt2[i] = p{i: g.i + er*999999, j: g.j + ec*999999}
	}

	return galaxiesPt1, galaxiesPt2
}

func abs(p1, p2 p) int {
	var di, dj int
	if p1.i > p2.i {
		di = p1.i - p2.i
	} else {
		di = p2.i - p1.i
	}

	if p1.j > p2.j {
		dj = p1.j - p2.j
	} else {
		dj = p2.j - p1.j
	}

	return di + dj
}
