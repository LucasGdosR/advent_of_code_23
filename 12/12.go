package main

import (
	"aoc/2023/common"
	"bufio"
	"strings"
)

var memo = make(map[string]int)

func main() {
	conditionRecords := common.Open("input")
	defer conditionRecords.Close()

	var results common.Results[int, int]
	scanner := bufio.NewScanner(conditionRecords)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		dotsNHashes, nums := line[0], line[1]
		results.Part1 += possibleArrangements(dotsNHashes, nums)
		results.Part2 += possibleArrangements(unfold(dotsNHashes, nums))
	}

	println(results.Part1, results.Part2)
}

func possibleArrangements(dnd, nums string) int {
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
		count = possibleArrangements(dnd[1:], nums)
	case '#':
		count = pound(dnd, nums)
	case '?':
		count = possibleArrangements(dnd[1:], nums) + pound(dnd, nums)
	}
	memo[key] = count
	return count
}

func pound(dnd, nums string) int {
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
			return possibleArrangements(dnd[num+1:], "")
		} else {
			return possibleArrangements(dnd[num+1:], nums[i+1:])
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
