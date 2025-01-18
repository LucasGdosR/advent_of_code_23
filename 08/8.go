package main

import (
	"aoc/2023/common"
	"bufio"
	"fmt"
)

type key struct {
	str string
	b   byte
}

func main() {
	file := common.Open("input")
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

	var part1Count int
	l := len(LR)
	part2Count := l
	for _, start := range startingPoints {
		count := 0
		for curr := start; curr[2] != 'Z'; count++ {
			curr = LRMap[key{str: curr, b: LR[count%l]}]
		}
		if start == "AAA" {
			part1Count = count
		}
		part2Count *= count / l
	}

	fmt.Println(part1Count, part2Count)
}
