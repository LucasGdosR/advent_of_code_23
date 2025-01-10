package main

import (
	"aoc/2023/common"
	"bufio"
	"fmt"
	"runtime"
)

const (
	GRID_SIZE           = 140
	GRID_PLUS_SENTINELS = GRID_SIZE + 2
	SENTINEL            = 1
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  findPartsAndGearRatiosST,
		MT_Impl:  findPartsAndGearRatiosMT,
		Part1Str: "Parts sum",
		Part2Str: "Gear ratios sum",
	}
	common.Benchmark(thisProgram, 1000)
}

func findPartsAndGearRatiosST() common.Results[int, int] {
	schematic := readSchematic()

	partNumbersSum := 0
	gearRatiosSum := 0
	parts := make([]int, 6)
	for i := SENTINEL; i < GRID_PLUS_SENTINELS-1-SENTINEL; i++ {
		for j := SENTINEL; j < GRID_PLUS_SENTINELS-1-SENTINEL; j++ {
			char := schematic[i][j]
			if char != '.' && (char < '0' || char > '9') {
				neighborCount := getNeighboringParts(schematic, i, j, parts)
				neighborSum := sumParts(parts, neighborCount)
				partNumbersSum += neighborSum
				if char == '*' && neighborCount == 2 {
					gearRatiosSum += parts[0] * parts[1]
				}
			}
		}
	}
	return common.Results[int, int]{Part1: partNumbersSum, Part2: gearRatiosSum}
}

func findPartsAndGearRatiosMT() common.Results[int, int] {
	schematic := readSchematic()

	numWorkers := runtime.GOMAXPROCS(0)
	partialResults := make(chan common.Results[int, int], numWorkers)
	linesPerWorker := GRID_SIZE / numWorkers

	for i := 0; i < numWorkers; i++ {
		start := i*linesPerWorker + SENTINEL
		end := start + linesPerWorker
		if i == numWorkers-1 {
			end = GRID_PLUS_SENTINELS - SENTINEL
		}
		go workerJob(start, end, partialResults, schematic)
	}
	var total common.Results[int, int]
	for i := 0; i < numWorkers; i++ {
		r := <-partialResults
		total.Part1 += r.Part1
		total.Part2 += r.Part2
	}
	close(partialResults)

	return total
}

func workerJob(start, end int, partialResults chan common.Results[int, int], schematic [][]byte) {
	parts := make([]int, 6)
	var workerResult common.Results[int, int]
	for i := start; i < end; i++ {
		for j := 1; j < GRID_PLUS_SENTINELS-2; j++ {
			char := schematic[i][j]
			if char != '.' && (char < '0' || char > '9') {
				neighborCount := getNeighboringParts(schematic, i, j, parts)
				workerResult.Part1 += sumParts(parts, neighborCount)
				if char == '*' && neighborCount == 2 {
					workerResult.Part2 += parts[0] * parts[1]
				}
			}
		}
	}
	partialResults <- workerResult
}

func readSchematic() [][]byte {
	schematic := common.Open("input")
	defer schematic.Close()

	matrix := make([][]byte, GRID_PLUS_SENTINELS)
	for i := range matrix {
		matrix[i] = make([]byte, GRID_PLUS_SENTINELS)
	}

	scanner := bufio.NewScanner(schematic)
	row := 1
	for scanner.Scan() {
		line := scanner.Text()
		copy(matrix[row][1:], []byte(line))
		row++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return matrix
}

func getNeighboringParts(m [][]byte, i int, j int, parts []int) int {
	neighborCount := 0
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			nextI := i + di
			nextJ := j + dj
			char := m[nextI][nextJ]
			if isNum(char) {
				start := searchLeft(m, nextI, nextJ)
				neighboringPart, end := readRight(m, nextI, start)
				dj = end - j
				parts[neighborCount] = neighboringPart
				neighborCount++
			}
		}
	}
	return neighborCount
}

func searchLeft(m [][]byte, i int, j int) int {
	char := m[i][j-1]
	if '0' <= char && char <= '9' {
		return searchLeft(m, i, j-1)
	} else {
		return j
	}
}

func readRight(m [][]byte, i int, j int) (int, int) {
	var thisPart int = 0
	for isNum(m[i][j]) {
		thisPart = thisPart*10 + int(m[i][j]-'0')
		j++
	}
	return thisPart, j - 1
}

func sumParts(parts []int, count int) int {
	sum := 0
	for i := 0; i < count; i++ {
		sum += parts[i]
	}
	return sum
}

func isNum(char byte) bool {
	return '0' <= char && char <= '9'
}
