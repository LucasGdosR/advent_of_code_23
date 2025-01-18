package main

import (
	"aoc/2023/common"
	"bufio"
)

type key struct {
	str string
	b   byte
}

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  countStepsToExitST,
		MT_Impl:  countStepsToExitMT,
		Part1Str: "Steps to exit",
		Part2Str: "Steps to exit",
	}
	common.Benchmark(thisProgram, 1000)
}

func countStepsToExitST() common.Results[int, int] {
	LR, LRMap, startingPoints := parseInput()

	l := len(LR)
	results := common.Results[int, int]{Part2: l}
	for _, start := range startingPoints {
		count := 0
		for curr := start; curr[2] != 'Z'; count++ {
			curr = LRMap[key{str: curr, b: LR[count%l]}]
		}
		if start == "AAA" {
			results.Part1 = count
		}
		results.Part2 *= count / l
	}

	return results
}

func countStepsToExitMT() common.Results[int, int] {
	// This could be done in parallel, but it's not worth the trouble
	LR, LRMap, startingPoints := parseInput()

	l, points := len(LR), len(startingPoints)
	results := common.Results[int, int]{Part2: l}
	counts := make(chan int, len(startingPoints))
	for _, start := range startingPoints {
		go func() {
			count := 0
			for curr := start; curr[2] != 'Z'; count++ {
				curr = LRMap[key{str: curr, b: LR[count%l]}]
			}
			if start == "AAA" {
				results.Part1 = count
			}
			counts <- count / l
		}()
	}
	for i := 0; i < points; i++ {
		results.Part2 *= <-counts
	}
	close(counts)

	return results
}

func parseInput() (string, map[key]string, []string) {
	file := common.Open("input")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	LR := scanner.Text()
	scanner.Scan()

	LRMap := make(map[key]string)
	startingPoints := make([]string, 0, 6)
	for scanner.Scan() {
		line := scanner.Text()
		curr := line[:3]
		LRMap[key{str: curr, b: 'L'}] = line[7:10]
		LRMap[key{str: curr, b: 'R'}] = line[12:15]
		if line[2] == 'A' {
			startingPoints = append(startingPoints, curr)
		}
	}
	return LR, LRMap, startingPoints
}
