package main

import (
	"aoc/2023/common"
	"math"
	"runtime"
	"sync"
	"syscall"
)

type p struct {
	i, j int
}

type task struct {
	tag, start, end int
}

const (
	UNIVERSE_LENGTH = 140
	GALAXY_COUNT    = 443
	LINE_BREAK      = 1
	TASK_READ_ROW   = iota
	TASK_READ_COL
	TASK_EXPAND
	TASK_SUM
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  measureUniverseDistancesST,
		MT_Impl:  measureUniverseDistancesMT,
		Part1Str: "Sum of distances between galaxies",
		Part2Str: "Sum of distances between galaxies",
	}
	common.Benchmark(thisProgram, 1000)
}

func measureUniverseDistancesST() common.Results[int, int] {
	var galaxiesPt1, galaxiesPt2 []p
	galaxiesPt1, galaxiesPt2 = expandUniverse(readImageST())
	return sumDistances(galaxiesPt1, galaxiesPt2, 0, GALAXY_COUNT)
}

func measureUniverseDistancesMT() common.Results[int, int] {
	mappedFile := common.Mmap("input")
	file := mappedFile.File
	defer syscall.Munmap(file)

	// Init workers variables
	numWorkers := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	taskChan := make(chan task, numWorkers)
	galaxiesChan := make(chan []p, numWorkers)
	distancesChan := make(chan common.Results[int, int], numWorkers)
	emptyRows := make([]int, UNIVERSE_LENGTH)
	emptyCols := make([]int, UNIVERSE_LENGTH)
	var galaxies, gp1, gp2 []p
	for i := range emptyRows {
		emptyRows[i] = 1
		emptyCols[i] = 1
	}

	for i := 0; i < numWorkers; i++ {
		go worker(&wg,
			taskChan,
			file,
			galaxiesChan,
			distancesChan,
			emptyRows,
			emptyCols,
			&galaxies,
			&gp1,
			&gp2)
	}

	// Read image
	linesPerWorker := (UNIVERSE_LENGTH + numWorkers - 1) / numWorkers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		start := i * linesPerWorker
		end := start + linesPerWorker
		if i == numWorkers-1 {
			end = UNIVERSE_LENGTH
		}
		taskChan <- task{tag: TASK_READ_ROW, start: start, end: end}
		taskChan <- task{tag: TASK_READ_COL, start: start, end: end}
	}

	// Merge galaxies
	galaxies = make([]p, 0, GALAXY_COUNT)
	gp1 = make([]p, GALAXY_COUNT)
	gp2 = make([]p, GALAXY_COUNT)
	for i := 0; i < numWorkers; i++ {
		galaxies = append(galaxies, <-galaxiesChan...)
	}
	close(galaxiesChan)
	wg.Wait()

	// Pre-compute emptyRows / Cols
	for erc, ecc, i := 0, 0, 0; i < UNIVERSE_LENGTH; i++ {
		erc += emptyRows[i]
		ecc += emptyCols[i]
		emptyRows[i] = erc
		emptyCols[i] = ecc
	}

	// Expand universe
	galaxiesPerWorker := (len(galaxies) + numWorkers - 1) / numWorkers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		start := i * galaxiesPerWorker
		end := start + galaxiesPerWorker
		if i == numWorkers-1 {
			end = len(galaxies)
		}
		taskChan <- task{tag: TASK_EXPAND, start: start, end: end}
	}
	wg.Wait()

	distanceCalculationCount := GALAXY_COUNT * (GALAXY_COUNT - 1) / 2
	step := (distanceCalculationCount + numWorkers - 1) / numWorkers
	for start, i := 0, 0; i < numWorkers; i++ {
		var end int
		if i == numWorkers-1 {
			end = GALAXY_COUNT
		} else {
			// This distributes work evenly between threads. Trust me.
			end = GALAXY_COUNT - int(math.Sqrt(float64(GALAXY_COUNT*GALAXY_COUNT-2*GALAXY_COUNT*(start+1)+start*start+2*start-2*step+1)))
		}
		taskChan <- task{tag: TASK_SUM, start: start, end: end}
		start = end
	}
	close(taskChan)

	var results common.Results[int, int]
	for i := 0; i < numWorkers; i++ {
		r := <-distancesChan
		results.Part1 += r.Part1
		results.Part2 += r.Part2
	}
	close(distancesChan)

	return results
}

func worker(wg *sync.WaitGroup,
	taskChan chan task,
	universe []byte,
	galaxiesChan chan []p,
	distancesChan chan common.Results[int, int],
	emptyRows, emptyCols []int,
	galaxies, gp1, gp2 *[]p) {

	for t := range taskChan {
		switch t.tag {
		case TASK_READ_ROW:
			partialGalaxies := make([]p, 0, GALAXY_COUNT/2)
			for i := t.start; i < t.end; i++ {
				for j := 0; j < UNIVERSE_LENGTH; j++ {
					if universe[(UNIVERSE_LENGTH+LINE_BREAK)*i+j] == '#' {
						partialGalaxies = append(partialGalaxies, p{i: i, j: j})
						emptyRows[i] = 0
					}
				}
			}
			galaxiesChan <- partialGalaxies

		case TASK_READ_COL:
			for j := t.start; j < t.end; j++ {
				for i := 0; i < UNIVERSE_LENGTH; i++ {
					if universe[(UNIVERSE_LENGTH+LINE_BREAK)*i+j] == '#' {
						emptyCols[j] = 0
						break
					}
				}
			}
			wg.Done()

		case TASK_EXPAND:
			fillSpace(*galaxies, *gp1, *gp2, emptyRows, emptyCols, t.start, t.end)
			wg.Done()

		case TASK_SUM:
			distancesChan <- sumDistances(*gp1, *gp2, t.start, t.end)
		}
	}
}

func readImageST() ([]int, []int, []p) {
	mappedFile := common.Mmap("input")
	universe := mappedFile.File
	defer syscall.Munmap(universe)

	emptyRows := make([]int, UNIVERSE_LENGTH)
	emptyCols := make([]int, UNIVERSE_LENGTH)
	for i := range emptyCols {
		emptyCols[i] = 1
	}
	galaxies := make([]p, 0, GALAXY_COUNT)

	for emptyRowCount, i := 0, 0; i < UNIVERSE_LENGTH; i++ {
		emptyRow := true
		for j := 0; j < UNIVERSE_LENGTH; j++ {
			if universe[(UNIVERSE_LENGTH+LINE_BREAK)*i+j] == '#' {
				galaxies = append(galaxies, p{i: i, j: j})
				emptyRow = false
				emptyCols[j] = 0
			}
		}
		if emptyRow {
			emptyRowCount++
		}
		emptyRows[i] = emptyRowCount
	}

	emptyColCount := 0
	for i, empty := range emptyCols {
		emptyColCount += empty
		emptyCols[i] = emptyColCount
	}

	return emptyRows, emptyCols, galaxies
}

func expandUniverse(emptyRows, emptyCols []int, galaxies []p) ([]p, []p) {
	galaxiesPt1 := make([]p, len(galaxies))
	galaxiesPt2 := make([]p, len(galaxies))

	fillSpace(galaxies, galaxiesPt1, galaxiesPt2, emptyRows, emptyCols, 0, GALAXY_COUNT)

	return galaxiesPt1, galaxiesPt2
}

func fillSpace(gs, g1, g2 []p, ers, ecs []int, start, end int) {
	for i := start; i < end; i++ {
		g := gs[i]
		er, ec := ers[g.i], ecs[g.j]
		g1[i] = p{i: g.i + er, j: g.j + ec}
		g2[i] = p{i: g.i + er*999999, j: g.j + ec*999999}
	}
}

func sumDistances(g1, g2 []p, start, end int) common.Results[int, int] {
	var results common.Results[int, int]
	for i, thisG := range g1[start:end] {
		for j, thatG := range g1[start+i+1:] {
			results.Part1 += abs(thisG, thatG)
			results.Part2 += abs(g2[start+i], g2[start+i+j+1])
		}
	}
	return results
}

func abs(p1, p2 p) int {
	var di, dj int
	if p1.i > p2.i {
		di = p1.i - p2.i
	} else {
		di = p2.i - p1.i
	}

	if p1.j > p2.j {
		dj = p1.j - p2.j
	} else {
		dj = p2.j - p1.j
	}

	return di + dj
}
