package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"runtime"
	"syscall"

	"github.com/dterei/gotsc"
)

type results struct {
	part1, part2 int
}

type benchmark struct {
	min, avg, max uint64
}

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
	printST = true
	printMT = true
)

func main() {
	const runs = 1000
	singleT := benchmark{min: math.MaxUint64, avg: 0, max: 0}
	multiT := benchmark{min: math.MaxUint64, avg: 0, max: 0}

	tsc := gotsc.TSCOverhead()
	fmt.Println("TSC Overhead:", tsc)

	for i := 0; i < runs; i++ {
		timeST := benchmarkSingleThreaded()
		singleT.min = min(singleT.min, timeST)
		singleT.avg += timeST / runs
		singleT.max = max(singleT.max, timeST)

		timeC := benchmarkConcurrent()
		multiT.min = min(multiT.min, timeC)
		multiT.avg += timeC / runs
		multiT.max = max(multiT.max, timeC)
	}

	fmt.Println("Single threaded (cycles):", "Min:", singleT.min, "Avg:", singleT.avg, "Max:", singleT.max)
	fmt.Println("Multi threaded  (cycles):", "Min:", multiT.min, "Avg:", multiT.avg, "Max:", multiT.max)
}

func benchmarkSingleThreaded() uint64 {
	start := gotsc.BenchStart()
	calibrationSums := fixDocumentST()
	end := gotsc.BenchEnd()

	if printST {
		fmt.Println("Calibration sum (part 1, ST):", calibrationSums.part1)
		fmt.Println("Calibration sum (part 2, ST):", calibrationSums.part2)
		printST = false
	}

	return end - start
}

func benchmarkConcurrent() uint64 {
	start := gotsc.BenchStart()
	calibrationSums := fixDocumentMT()
	end := gotsc.BenchEnd()

	if printMT {
		fmt.Println("Calibration sum (part 1, MT):", calibrationSums.part1)
		fmt.Println("Calibration sum (part 2, MT):", calibrationSums.part2)
		printMT = false
	}

	return end - start
}

func fixDocumentST() results {
	calibrationDocument, err := os.Open("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer calibrationDocument.Close()

	scanner := bufio.NewScanner(calibrationDocument)
	var calibrationSums results

	var buf = make([]byte, 64)
	for scanner.Scan() {
		line := scanner.Text()

		first_digit := findFirstDigit(line)
		last_digit := findLastDigit(line)
		calibrationSums.part1 += int(first_digit*10 + last_digit)

		firstNumber := findFirstNumber(line)
		lastNumber := findLastNumber(line, buf)
		calibrationSums.part2 += int(firstNumber*10 + lastNumber)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return calibrationSums
}

func findFirstDigit(s string) byte {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return s[i] - '0'
		}
	}
	return 0
}

func findLastDigit(s string) byte {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			return s[i] - '0'
		}
	}
	return 0
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

func fixDocumentMT() results {
	file, err := os.OpenFile("input", os.O_RDONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return results{}
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file size: %v\n", err)
		return results{}
	}
	size := fi.Size()

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error mapping file: %v\n", err)
		return results{}
	}
	defer syscall.Munmap(data)

	numWorkers := runtime.GOMAXPROCS(0)
	partialResults := make(chan results, numWorkers)
	bytesPerWorker := size / int64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := int64(i) * bytesPerWorker
		end := start + bytesPerWorker
		if i == numWorkers-1 {
			end = size
		}

		go func(start, end int64) {
			if start > 0 {
				for data[start-1] != '\n' {
					start--
				}
			}
			if end < int64(len(data)) {
				for data[end-1] != '\n' {
					end--
				}
			}
			partial := processMemRange(data, start, end)
			partialResults <- partial
		}(start, end)
	}

	var total results
	for i := 0; i < numWorkers; i++ {
		r := <-partialResults
		total.part1 += r.part1
		total.part2 += r.part2
	}

	return total
}

func processMemRange(data []byte, start int64, end int64) results {
	var partial results

	var buf = make([]byte, 64)
	for start < end {
		lineEnd := start
		for lineEnd < end && data[lineEnd] != '\n' {
			lineEnd++
		}

		line := string(data[start:lineEnd])

		first_digit := findFirstDigit(line)
		last_digit := findLastDigit(line)
		partial.part1 += int(first_digit*10 + last_digit)

		firstNumber := findFirstNumber(line)
		lastNumber := findLastNumber(line, buf)
		partial.part2 += int(firstNumber*10 + lastNumber)

		start = lineEnd + 1
	}

	return partial
}
