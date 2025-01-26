package main

import (
	"aoc/2023/common"
	"bufio"
	"bytes"
	"runtime"
	"syscall"
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  findAllMirrorsST,
		MT_Impl:  findAllMirrorsMT,
		Part1Str: "Summarized notes",
		Part2Str: "Summarized notes",
	}
	common.Benchmark(thisProgram, 1000)
}

func findAllMirrorsST() common.Results[int, int] {
	file := common.Open("input")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var results common.Results[int, int]
	pattern := make([][]byte, 32)
	i := 0
	for ; scanner.Scan(); i++ {
		line := scanner.Bytes()
		if len(line) == 0 {
			v := findMirrors(pattern[:i], -1)
			results.Part1 += v
			results.Part2 += findSmudge(pattern[:i], v)
			i = 0
			scanner.Scan()
			line = scanner.Bytes()
		}
		lineCopy := make([]byte, len(line))
		copy(lineCopy, line)
		pattern[i] = lineCopy
	}
	// Solve the last pattern
	v := findMirrors(pattern[:i], -1)
	results.Part1 += v
	results.Part2 += findSmudge(pattern[:i], v)

	return results
}

func findAllMirrorsMT() common.Results[int, int] {
	mappedFile := common.Mmap("input")
	file := mappedFile.File
	size := mappedFile.Size
	defer syscall.Munmap(file)

	numWorkers := runtime.GOMAXPROCS(0)
	bytesPerWorker := (int(size) + numWorkers) / numWorkers
	partialResults := make(chan common.Results[int, int], numWorkers)
	for i := 0; i < numWorkers; i++ {
		start := i * bytesPerWorker
		end := start + bytesPerWorker
		if i == numWorkers-1 {
			end = int(size) + 2 // Two missing '\n's in the last pattern
		}
		go func(s, e int) {
			if start > 0 {
				for file[start-1] != '\n' || file[start-2] != '\n' {
					start--
				}
			}
			if end < int(size) {
				for file[end-1] != '\n' || file[end-2] != '\n' {
					end--
				}
			}

			partialResults <- parseAndSolve(file[start : end-2]) // Skip 2 last '\n'
		}(start, end)
	}
	var results common.Results[int, int]
	for i := 0; i < numWorkers; i++ {
		r := <-partialResults
		results.Part1 += r.Part1
		results.Part2 += r.Part2
	}
	close(partialResults)

	return results
}

func parseAndSolve(patterns []byte) common.Results[int, int] {
	pattern := make([][]byte, 32)
	end := len(patterns)
	row := make([]byte, 0, 24)
	var results common.Results[int, int]

	j := 0
	// Read lines
	for i := 0; i < end; i++ {
		b := patterns[i]
		// Append row to pattern
		if b == '\n' {
			pattern[j] = row
			row = make([]byte, 0, 24)
			j++
			// Solve pattern, advance cursor, reuse pattern memory
			if patterns[i+1] == '\n' {
				v := findMirrors(pattern[:j], -1)
				results.Part1 += v
				results.Part2 += findSmudge(pattern[:j], v)
				j = 0
				i++
			}
			// Append char to row
		} else {
			row = append(row, b)
		}
	}
	// Solve the last pattern
	pattern[j] = row
	j++
	v := findMirrors(pattern[:j], -1)
	results.Part1 += v
	results.Part2 += findSmudge(pattern[:j], v)

	return results
}

func findMirrors(p [][]byte, incorrectV int) int {
	// Scan first line looking for the same char in a row.
	row := p[0]
	for i, c := range row[1:] {
		if i+1 == incorrectV {
			continue
		}
		if row[i] == byte(c) {
			// If found, see if it repeats in the other rows.
			isMatch := true
			for _, r := range p[1:] {
				if r[i] != r[i+1] {
					isMatch = false
					break
				}
			}
			if isMatch {
				// If it does, check the whole block.
				if assertVerticalMirrorring(p, i+1) {
					return i + 1
				}
			}
		}
	}
	// If no matches, look for columns.
	for i := 1; i < len(p); i++ {
		if 100*i == incorrectV {
			continue
		}
		if p[i-1][0] == p[i][0] {
			if bytes.Equal(p[i-1], p[i]) {
				if assertHorizontalMirrorring(p, i) {
					return 100 * i
				}
			}
		}
	}
	return 0
}

func findSmudge(p [][]byte, incorrectV int) int {
	for _, row := range p {
		for i, c := range row {
			temp := c
			if c == '.' {
				row[i] = '#'
			} else {
				row[i] = '.'
			}
			v := findMirrors(p, incorrectV)
			if v != 0 && v != incorrectV {
				return v
			} else {
				row[i] = temp
			}
		}
	}
	return 0
}

func assertHorizontalMirrorring(p [][]byte, i int) bool {
	for u, d := i-2, i+1; u >= 0 && d < len(p); u, d = u-1, d+1 {
		if !bytes.Equal(p[u], p[d]) {
			return false
		}
	}
	return true
}

func assertVerticalMirrorring(p [][]byte, i int) bool {
	for _, row := range p {
		for l, r := i-2, i+1; l >= 0 && r < len(p[0]); l, r = l-1, r+1 {
			if row[l] != row[r] {
				return false
			}
		}
	}
	return true
}
