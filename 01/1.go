package main

import (
	"aoc/2023/common"
	"bufio"
	"fmt"
	"regexp"
)

var (
	forwardRe, _ = regexp.Compile(`(one|two|three|four|five|six|seven|eight|nine|\d)`)
	reverseRe, _ = regexp.Compile(`(eno|owt|eerht|ruof|evif|xis|neves|thgie|enin|\d)`)
	numberMap    = map[string]byte{
		"one": 1, "1": 1,
		"two": 2, "2": 2,
		"three": 3, "3": 3,
		"four": 4, "4": 4,
		"five": 5, "5": 5,
		"six": 6, "6": 6,
		"seven": 7, "7": 7,
		"eight": 8, "8": 8,
		"nine": 9, "9": 9,
	}
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl: fixDocumentST,
		MT_Impl: func() common.Results[int, int] {
			return common.SolveCommonCaseMmapLinesInt(fixDocumentRange)
		},
		Part1Str: "Calibration sum",
		Part2Str: "Calibration sum",
	}
	common.Benchmark(thisProgram, 1000)
}

func fixDocumentST() common.Results[int, int] {
	calibrationDocument := common.Open("input")
	defer calibrationDocument.Close()

	scanner := bufio.NewScanner(calibrationDocument)
	var calibrationSums common.Results[int, int]

	var buf = make([]byte, 64)
	for scanner.Scan() {
		line := scanner.Text()
		part1, part2 := calibrateLine(line, buf)
		calibrationSums.Part1 += part1
		calibrationSums.Part2 += part2
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return calibrationSums
}

func fixDocumentRange(data []byte, start int64, end int64) common.Results[int, int] {
	var partial common.Results[int, int]

	var buf = make([]byte, 64)
	for start < end {
		lineEnd := start
		for lineEnd < end && data[lineEnd] != '\n' {
			lineEnd++
		}

		line := string(data[start:lineEnd])
		part1, part2 := calibrateLine(line, buf)
		partial.Part1 += part1
		partial.Part2 += part2

		start = lineEnd + 1
	}

	return partial
}

func calibrateLine(line string, buf []byte) (int, int) {
	first_digit := findFirstDigit(line)
	last_digit := findLastDigit(line)

	firstNumber := findFirstNumber(line)
	lastNumber := findLastNumber(line, buf)

	return int(first_digit*10 + last_digit), int(firstNumber*10 + lastNumber)
}

func findFirstDigit(s string) byte {
	var digit byte = 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			digit = s[i] - '0'
			break
		}
	}
	return digit
}

func findLastDigit(s string) byte {
	var digit byte = 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			digit = s[i] - '0'
			break
		}
	}
	return digit
}

func findFirstNumber(s string) byte {
	match := forwardRe.FindString(s)
	return numberMap[match]
}

func findLastNumber(s string, buf []byte) byte {
	match := reverseRe.FindString(reverse(s, buf))
	return numberMap[reverse(match, buf)]
}

func reverse(s string, buf []byte) string {
	size := len(s)
	buf = buf[:size]
	for i := 0; i < size; i++ {
		buf[size-1-i] = s[i]
	}
	return string(buf)
}
