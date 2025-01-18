package main

import (
	"aoc/2023/common"
	"bufio"
	"strings"
)

const SEQ_LEN = 21

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  extrapolateST,
		MT_Impl:  extrapolateMT,
		Part1Str: "Sum of next terms",
		Part2Str: "Sum of previous terms",
	}
	common.Benchmark(thisProgram, 1000)
}

func extrapolateST() common.Results[int, int] {
	file := common.Open("input")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	triangle := make([][]int, SEQ_LEN)
	for i := range triangle {
		triangle[i] = make([]int, SEQ_LEN+2-i)
	}

	var results common.Results[int, int]
	for scanner.Scan() {
		parseSeq(scanner.Text(), triangle)
		depth := fillTriangle(triangle)
		results.Part1 += getNext(triangle, depth)
		results.Part2 += getPrev(triangle, depth)
	}
	return results
}

func extrapolateMT() common.Results[int, int] {
	return common.SolveCommonCaseMmapLinesInt(extrapolateRange)
}

func extrapolateRange(seqs []byte, start, end int64) common.Results[int, int] {
	var results common.Results[int, int]

	triangle := make([][]int, SEQ_LEN)
	for i := range triangle {
		triangle[i] = make([]int, SEQ_LEN+2-i)
	}

	for start < end {
		lineEnd := start
		for lineEnd < end && seqs[lineEnd] != '\n' {
			lineEnd++
		}

		parseSeq(string(seqs[start:lineEnd]), triangle)
		depth := fillTriangle(triangle)
		results.Part1 += getNext(triangle, depth)
		results.Part2 += getPrev(triangle, depth)

		start = lineEnd + 1
	}

	return results
}

func parseSeq(line string, triangle [][]int) {
	seq := strings.Fields(line)
	for i, num := range seq {
		triangle[0][i+1] = common.Atoi(num)
	}
}

func fillTriangle(t [][]int) int {
	depth := SEQ_LEN - 1
	for i := 1; i < SEQ_LEN; i++ {
		allZero := true
		nextRowLen := SEQ_LEN - i + 1
		for j := 1; j < nextRowLen; j++ {
			v := t[i-1][j+1] - t[i-1][j]
			t[i][j] = v
			allZero = allZero && v == 0
		}
		if allZero {
			depth = i - 1
			break
		}
	}
	return depth
}

func getNext(t [][]int, depth int) int {
	t[depth][SEQ_LEN-depth+1] = t[depth][SEQ_LEN-depth]
	for i := depth - 1; i >= 0; i-- {
		t[i][SEQ_LEN-i+1] = t[i+1][SEQ_LEN-i] + t[i][SEQ_LEN-i]
	}
	return t[0][SEQ_LEN+1]
}

func getPrev(t [][]int, depth int) int {
	t[depth][0] = t[depth][1]
	for i := depth - 1; i >= 0; i-- {
		t[i][0] = t[i][1] - t[i+1][0]
	}
	return t[0][0]
}
