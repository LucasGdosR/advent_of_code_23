package common

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"syscall"

	"github.com/dterei/gotsc"
)

type benchmark struct {
	min, avg, max uint64
}

type Results[T1, T2 any] struct {
	Part1 T1
	Part2 T2
}

type Benchmarkee[T1, T2 any] struct {
	ST_Impl, MT_Impl   func() Results[T1, T2]
	Part1Str, Part2Str string
}

type MappedFile struct {
	File []byte
	Size int64
}

var printST, printMT = true, true

func Benchmark[T1, T2 any](benchmarkee Benchmarkee[T1, T2], runs int) {
	singleT := benchmark{min: math.MaxUint64, avg: 0, max: 0}
	multiT := benchmark{min: math.MaxUint64, avg: 0, max: 0}

	runsDivider := uint64(runs)
	for i := 0; i < runs; i++ {
		timeST := benchmarkSingleThreaded(benchmarkee)
		singleT.min = min(singleT.min, timeST)
		singleT.avg += timeST / runsDivider
		singleT.max = max(singleT.max, timeST)

		timeC := benchmarkMultiThreaded(benchmarkee)
		multiT.min = min(multiT.min, timeC)
		multiT.avg += timeC / runsDivider
		multiT.max = max(multiT.max, timeC)
	}

	fmt.Println("Single threaded (cycles):", "Min:", singleT.min, "Avg:", singleT.avg, "Max:", singleT.max)
	fmt.Println("Multi threaded  (cycles):", "Min:", multiT.min, "Avg:", multiT.avg, "Max:", multiT.max)
}

func benchmarkSingleThreaded[T1, T2 any](benchmarkee Benchmarkee[T1, T2]) uint64 {
	start := gotsc.BenchStart()
	results := benchmarkee.ST_Impl()
	end := gotsc.BenchEnd()

	if printST {
		fmt.Println(benchmarkee.Part1Str, "(part 1, ST):", results.Part1)
		fmt.Println(benchmarkee.Part2Str, "(part 2, ST):", results.Part2)
		printST = false
	}

	return end - start
}

func benchmarkMultiThreaded[T1, T2 any](benchmarkee Benchmarkee[T1, T2]) uint64 {
	start := gotsc.BenchStart()
	results := benchmarkee.MT_Impl()
	end := gotsc.BenchEnd()

	if printMT {
		fmt.Println(benchmarkee.Part1Str, "(part 1, MT):", results.Part1)
		fmt.Println(benchmarkee.Part2Str, "(part 2, MT):", results.Part2)
		printMT = false
	}

	return end - start
}

func Open(filepath string) *os.File {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	return file
}

func Mmap(filepath string) *MappedFile {
	file := Open(filepath)
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting file size: %v\n", err)
		os.Exit(1)
	}
	size := fi.Size()

	data, err := syscall.Mmap(int(file.Fd()), 0, int(size),
		syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error mapping file: %v\n", err)
		os.Exit(1)
	}

	return &MappedFile{File: data, Size: size}
}

func SolveCommonCaseMmapLinesInt(processMemRange func([]byte, int64, int64) Results[int, int]) Results[int, int] {
	mappedFile := Mmap("input")
	defer syscall.Munmap(mappedFile.File)

	numWorkers := runtime.GOMAXPROCS(0)
	partialResults := make(chan Results[int, int], numWorkers)

	MmapBacktrackingLinesSolution(mappedFile, partialResults, processMemRange, numWorkers)

	var total Results[int, int]
	for i := 0; i < numWorkers; i++ {
		r := <-partialResults
		total.Part1 += r.Part1
		total.Part2 += r.Part2
	}
	close(partialResults)

	return total
}

func MmapBacktrackingLinesSolution[T1, T2 any](
	mappedFile *MappedFile,
	partialResults chan Results[T1, T2],
	processMemRange func([]byte, int64, int64) Results[T1, T2],
	numWorkers int) {

	file := mappedFile.File
	size := mappedFile.Size
	bytesPerWorker := size / int64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := int64(i) * bytesPerWorker
		end := start + bytesPerWorker
		if i == numWorkers-1 {
			end = size
		}

		go func(start, end int64) {
			if start > 0 {
				for file[start-1] != '\n' {
					start--
				}
			}
			if end < int64(len(file)) {
				for file[end-1] != '\n' {
					end--
				}
			}
			partialResults <- processMemRange(file, start, end)
		}(start, end)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing int: %v\n", err)
		os.Exit(1)
	}
	return i
}
