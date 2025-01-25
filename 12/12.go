package main

import (
	"aoc/2023/common"
	"bufio"
	"strings"
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  countPossibleArragementsST,
		MT_Impl:  countPossibleArragementsMT,
		Part1Str: "Possible arrangements",
		Part2Str: "Possible arrangements",
	}
	common.Benchmark(thisProgram, 100)
}

func countPossibleArragementsST() common.Results[int, int] {
	conditionRecords := common.Open("input")
	defer conditionRecords.Close()

	memo := make(map[string]int, 260417)
	var results common.Results[int, int]
	scanner := bufio.NewScanner(conditionRecords)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		dotsNHashes, nums := line[0], line[1]
		results.Part1 += possibleArrangements(dotsNHashes, nums, memo)
		UDND, UNums := unfold(dotsNHashes, nums)
		results.Part2 += possibleArrangements(UDND, UNums, memo)
	}

	return results
}

func countPossibleArragementsMT() common.Results[int, int] {
	return common.SolveCommonCaseMmapLinesInt(countRangeArrangements)
}

func countRangeArrangements(data []byte, start int64, end int64) common.Results[int, int] {
	var partial common.Results[int, int]
	memo := make(map[string]int, 70000)
	for start < end {
		whitespace := start + 10
		for data[whitespace] != ' ' {
			whitespace++
		}
		dotsNHashes := string(data[start:whitespace])

		lineEnd := whitespace + 4
		for lineEnd < end && data[lineEnd] != '\n' {
			lineEnd++
		}
		nums := (string(data[whitespace+1 : lineEnd]))

		partial.Part1 += possibleArrangements(dotsNHashes, nums, memo)
		UDND, UNums := unfold(dotsNHashes, nums)
		partial.Part2 += possibleArrangements(UDND, UNums, memo)

		start = lineEnd + 1
	}

	return partial
}

func possibleArrangements(dnd, nums string, memo map[string]int) int {
	key := dnd + nums
	if v, ok := memo[key]; ok {
		return v
	}

	if len(nums) == 0 {
		if strings.Contains(dnd, "#") {
			return 0
		}
		return 1
	}

	if len(dnd) == 0 {
		return 0
	}

	var count int
	switch dnd[0] {
	case '.':
		count = possibleArrangements(dnd[1:], nums, memo)
	case '#':
		count = pound(dnd, nums, memo)
	case '?':
		count = possibleArrangements(dnd[1:], nums, memo) + pound(dnd, nums, memo)
	}
	memo[key] = count
	return count
}

func pound(dnd, nums string, memo map[string]int) int {
	var num int
	i := strings.IndexByte(nums, ',')
	if i == -1 {
		num = common.Atoi(nums)
	} else {
		num = common.Atoi(nums[:i])
	}
	if num > len(dnd) || strings.Contains(dnd[:num], ".") {
		return 0
	} else if num == len(dnd) {
		if i == -1 {
			return 1
		} else {
			return 0
		}
	} else if dnd[num] == '#' {
		return 0
	} else {
		if i == -1 {
			return possibleArrangements(dnd[num+1:], "", memo)
		} else {
			return possibleArrangements(dnd[num+1:], nums[i+1:], memo)
		}
	}
}

func unfold(dnd, nums string) (string, string) {
	unfoldedS := make([]string, 0, 5)
	unfoldedNums := make([]string, 0, 5)
	for i := 0; i < 5; i++ {
		unfoldedS = append(unfoldedS, dnd)
		unfoldedNums = append(unfoldedNums, nums)
	}
	return strings.Join(unfoldedS, "?"), strings.Join(unfoldedNums, ",")
}
