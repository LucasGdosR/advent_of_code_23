package main

import (
	"aoc/2023/common"
	"fmt"
	"os"
	"runtime"
	"sync"
)

type hash struct {
	loadN, loadW int
}

type task struct {
	tag, start, end int16
}

const (
	P_LEN      = 100
	LINE_BREAK = 1
	CALC_LOAD  = iota
	N
	W
	S
	E
)

var (
	platform   []byte
	indexCache map[uint64]int
	loadCache  map[int]uint64
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  measureLoadST,
		MT_Impl:  measureLoadMT,
		Part1Str: "North load",
		Part2Str: "North load after 1 billion cycles",
	}
	common.Benchmark(thisProgram, 1000)
}

func measureLoadST() common.Results[int, int] {
	rollingCache := make([]int16, P_LEN)
	indexCache = make(map[uint64]int)
	loadCache = make(map[int]uint64)

	var err error
	platform, err = os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var results common.Results[int, int]
	results.Part1 = firstIteration(rollingCache, 0, P_LEN)

	var firstOccurrence, secondOccurrence int = findCycle(rollingCache, 0, P_LEN)
	period := secondOccurrence - firstOccurrence
	toGo := (1000000000 - 1 - secondOccurrence) % period
	results.Part2 = int(loadCache[firstOccurrence+toGo] & 0xFFFFFFFF)

	return results
}

func measureLoadMT() common.Results[int, int] {
	// Init caches
	indexCache = make(map[uint64]int)
	loadCache = make(map[int]uint64)
	var results common.Results[int, int]

	// Read file
	var err error
	platform, err = os.ReadFile("input")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Init worker pool
	var wg sync.WaitGroup
	numWorkers := int16(runtime.GOMAXPROCS(0))
	taskChan := make(chan task, numWorkers)
	resultsChan := make(chan uint64, numWorkers)
	linesPerWorker := (P_LEN + numWorkers - 1) / numWorkers
	for i := int16(0); i < numWorkers; i++ {
		go worker(&wg, taskChan, resultsChan)
	}

	// Generic broadcasting func
	broadcastTask := func(tag int16) {
		for i := int16(0); i < numWorkers; i++ {
			start := i * linesPerWorker
			end := start + linesPerWorker
			if i == numWorkers {
				end = P_LEN
			}
			taskChan <- task{tag: tag, start: start, end: end}
		}
	}

	// firstIterationMT
	wg.Add(int(numWorkers))
	broadcastTask(N)
	wg.Wait()
	broadcastTask(CALC_LOAD)
	for i := int16(0); i < numWorkers; i++ {
		results.Part1 += int(<-resultsChan)
	}
	wg.Add(int(numWorkers))
	broadcastTask(W)
	wg.Wait()
	wg.Add(int(numWorkers))
	broadcastTask(S)
	wg.Wait()
	broadcastTask(E)
	var h uint64
	for i := int16(0); i < numWorkers; i++ {
		h += uint64(<-resultsChan)
	}
	indexCache[h] = 0

	// findCycleMT
	var firstOccurrence, secondOccurrence int
	i := 1
	for {
		wg.Add(int(numWorkers))
		broadcastTask(N)
		wg.Wait()
		wg.Add(int(numWorkers))
		broadcastTask(W)
		wg.Wait()
		wg.Add(int(numWorkers))
		broadcastTask(S)
		wg.Wait()
		broadcastTask(E)
		var _h uint64
		for i := int16(0); i < numWorkers; i++ {
			_h += uint64(<-resultsChan)
		}
		if j, ok := indexCache[_h]; ok {
			close(taskChan)
			close(resultsChan)
			firstOccurrence = j
			secondOccurrence = i
			break
		} else {
			indexCache[_h] = i
			loadCache[i] = _h
			i++
		}
	}
	period := secondOccurrence - firstOccurrence
	toGo := (1000000000 - 1 - secondOccurrence) % period
	results.Part2 = int(loadCache[firstOccurrence+toGo] & 0xFFFFFFFF)

	return results
}

func worker(wg *sync.WaitGroup, taskChan chan task, resultsChan chan uint64) {
	myRollingCache := make([]int16, P_LEN)
	for t := range taskChan {
		switch t.tag {
		case CALC_LOAD:
			resultsChan <- uint64(getLoad(t.start, t.end))
		case N:
			rollNorth(myRollingCache, t.start, t.end)
			wg.Done()
		case W:
			rollWest(myRollingCache, t.start, t.end)
			wg.Done()
		case S:
			rollSouth(myRollingCache, t.start, t.end)
			wg.Done()
		case E:
			hash := rollEast(myRollingCache, t.start, t.end)
			resultsChan <- uint64(hash.loadN) | uint64(hash.loadW<<32)
		}
	}
}

func firstIteration(rollingCache []int16, start, end int16) int {
	rollNorth(rollingCache, start, end)
	res := getLoad(0, P_LEN)
	rollWest(rollingCache, start, end)
	rollSouth(rollingCache, start, end)
	hash := rollEast(rollingCache, start, end)
	h := uint64(hash.loadN) | uint64(hash.loadW<<32)
	indexCache[h] = 0
	return res
}

func findCycle(rollingCache []int16, start, end int16) (int, int) {
	var firstOccurrence, secondOccurrence int
	i := 1
	for {
		rollNorth(rollingCache, start, end)
		rollWest(rollingCache, start, end)
		rollSouth(rollingCache, start, end)
		_hash := rollEast(rollingCache, start, end)
		_h := uint64(_hash.loadN) | uint64(_hash.loadW<<32)
		if j, ok := indexCache[_h]; ok {
			firstOccurrence = j
			secondOccurrence = i
			break
		} else {
			indexCache[_h] = i
			loadCache[i] = _h
			i++
		}
	}
	return firstOccurrence, secondOccurrence
}

func rollNorth(rollingCache []int16, start, end int16) {
	for i := start; i < end; i++ {
		rollingCache[i] = 0
	}
	for i := int16(0); i < P_LEN; i++ {
		for j := start; j < end; j++ {
			switch get(platform, i, j) {
			case 'O':
				set(platform, i, j, '.')
				set(platform, rollingCache[j], j, 'O')
				rollingCache[j]++
			case '#':
				rollingCache[j] = i + 1
			}
		}
	}
}

func rollWest(rollingCache []int16, start, end int16) {
	for i := start; i < end; i++ {
		rollingCache[i] = 0
	}
	for i := start; i < end; i++ {
		for j := int16(0); j < P_LEN; j++ {
			switch get(platform, i, j) {
			case 'O':
				set(platform, i, j, '.')
				set(platform, i, rollingCache[i], 'O')
				rollingCache[i]++
			case '#':
				rollingCache[i] = j + 1
			}
		}
	}
}

func rollSouth(rollingCache []int16, start, end int16) {
	for i := start; i < end; i++ {
		rollingCache[i] = P_LEN - 1
	}
	for i := int16(P_LEN - 1); i >= 0; i-- {
		for j := start; j < end; j++ {
			switch get(platform, i, j) {
			case 'O':
				set(platform, i, j, '.')
				set(platform, rollingCache[j], j, 'O')
				rollingCache[j]--
			case '#':
				rollingCache[j] = i - 1
			}
		}
	}
}

func rollEast(rollingCache []int16, start, end int16) hash {
	var h hash
	for i := start; i < end; i++ {
		rollingCache[i] = P_LEN - 1
	}
	for i := start; i < end; i++ {
		for j := int16(P_LEN - 1); j >= 0; j-- {
			switch get(platform, i, j) {
			case 'O':
				set(platform, i, j, '.')
				set(platform, i, rollingCache[i], 'O')
				h.loadN += P_LEN - int(i)
				h.loadW += int(rollingCache[i] ^ i)
				rollingCache[i]--
			case '#':
				rollingCache[i] = j - 1
			}
		}
	}
	return h
}

func getLoad(start, end int16) int {
	var load int
	for i := start; i < end; i++ {
		for j := int16(0); j < P_LEN; j++ {
			if get(platform, i, j) == 'O' {
				load += P_LEN - int(i)
			}
		}
	}
	return load
}

func get(platform []byte, i, j int16) byte {
	return platform[i*(P_LEN+LINE_BREAK)+j]
}

func set(platform []byte, i, j int16, b byte) {
	platform[i*(P_LEN+LINE_BREAK)+j] = b
}
