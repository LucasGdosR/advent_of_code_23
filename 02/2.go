package main

import (
	"aoc/2023/common"
	"bufio"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

var re, _ = regexp.Compile(`Game (\d+): (.*)`)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  solveColoredBallsST,
		MT_Impl:  solveColoredBallsMT,
		Part1Str: "IDs sum",
		Part2Str: "Power sum",
	}
	common.Benchmark(thisProgram, 1000)
}

func solveColoredBallsST() common.Results[int, int] {
	file := common.Open("input")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var idAndPowerSums common.Results[int, int]

	for scanner.Scan() {
		game := re.FindStringSubmatch(scanner.Text())
		id, power := solveGame(game)
		idAndPowerSums.Part1 += id
		idAndPowerSums.Part2 += power
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return idAndPowerSums
}

func solveColoredBallsMT() common.Results[int, int] {
	mappedFile := common.Mmap("input")
	defer syscall.Munmap(mappedFile.File)

	numWorkers := runtime.GOMAXPROCS(0)
	partialResults := make(chan common.Results[int, int], numWorkers)

	common.MmapBacktrackingLinesSolution(mappedFile, partialResults, solveColoredBallsRange, numWorkers)

	var total common.Results[int, int]
	for i := 0; i < numWorkers; i++ {
		r := <-partialResults
		total.Part1 += r.Part1
		total.Part2 += r.Part2
	}

	return total
}

func solveColoredBallsRange(data []byte, start int64, end int64) common.Results[int, int] {
	var partial common.Results[int, int]

	for start < end {
		lineEnd := start
		for lineEnd < end && data[lineEnd] != '\n' {
			lineEnd++
		}

		line := string(data[start:lineEnd])
		game := re.FindStringSubmatch(line)
		part1, part2 := solveGame(game)
		partial.Part1 += part1
		partial.Part2 += part2

		start = lineEnd + 1
	}

	return partial
}

func solveGame(game []string) (int, int) {
	id, _ := strconv.Atoi(game[1])
	var rMax, gMax, bMax = 0, 0, 0

	countColors := strings.FieldsFunc(game[2], func(r rune) bool {
		return r == ',' || r == ';'
	})

	for _, countColor := range countColors {
		cubes := strings.Fields(countColor)
		count, _ := strconv.Atoi(cubes[0])
		switch cubes[1] {
		case "red":
			if count > 12 {
				id = 0
			}
			rMax = max(rMax, count)
		case "green":
			if count > 13 {
				id = 0
			}
			gMax = max(gMax, count)
		case "blue":
			if count > 14 {
				id = 0
			}
			bMax = max(bMax, count)
		}
	}
	return id, rMax * gMax * bMax
}
