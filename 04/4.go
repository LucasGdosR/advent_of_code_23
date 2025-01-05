package main

import (
	"aoc/2023/common"
	"runtime"
	"syscall"
)

const (
	LINE_LENGTH                = 117
	CARD_COUNT                 = 211
	WINNERS_START, WINNERS_END = 10, 40
	WINNERS_LENGTH             = WINNERS_END - WINNERS_START
	NUMBERS_START, NUMBERS_END = 42, 116
	NUMBERS_LENGTH             = NUMBERS_END - NUMBERS_START
)

func main() {
	thisProgram := common.Benchmarkee[int, int]{
		ST_Impl:  sumPointsAndCountScratchCardsST,
		MT_Impl:  sumPointsAndCountScratchCardsMT,
		Part1Str: "Points sum",
		Part2Str: "Scratchcard count",
	}
	common.Benchmark(thisProgram, 1000)
}

func sumPointsAndCountScratchCardsST() common.Results[int, int] {
	mappedFile := common.Mmap("input")
	scratchcards := mappedFile.File
	defer syscall.Munmap(scratchcards)
	size := int(mappedFile.Size)

	var points = 0
	scratchcardsWinnings := make([]int, CARD_COUNT)
	totalScratchcards := 0
	for i, card := size-LINE_LENGTH+1, byte(CARD_COUNT-1); i >= 0; i, card = i-LINE_LENGTH, card-1 {
		winners := scratchcards[i+WINNERS_START : i+WINNERS_END]
		numbers := scratchcards[i+NUMBERS_START : i+NUMBERS_END]

		winnerSet := makeWinnersSet(winners)
		matches := getMatches(winnerSet, numbers)
		points += 1 << (matches - 1)

		scratchcardsWon := 1
		for i := byte(1); i <= matches; i++ {
			scratchcardsWon += scratchcardsWinnings[card+i]
		}
		scratchcardsWinnings[card] = scratchcardsWon
		totalScratchcards += scratchcardsWon
	}

	return common.Results[int, int]{Part1: points, Part2: totalScratchcards}
}

func sumPointsAndCountScratchCardsMT() common.Results[int, int] {
	mappedFile := common.Mmap("input")
	scratchcards := mappedFile.File
	defer syscall.Munmap(scratchcards)

	numWorkers := runtime.GOMAXPROCS(0)
	linesPerWoker := (CARD_COUNT + numWorkers - 1) / numWorkers

	partialPoints := make(chan int, numWorkers)
	matches := make([][]int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := i * linesPerWoker
		end := start + linesPerWoker
		if i == numWorkers-1 {
			end = CARD_COUNT
		}
		partialMatches := make([]int, end-start)
		matches[i] = partialMatches

		go func(start, end int, partialMatches []int) {
			myPoints := 0
			j := 0
			for i, card := start*LINE_LENGTH, start; card < end; i, card = i+LINE_LENGTH, card+1 {
				winners := scratchcards[i+WINNERS_START : i+WINNERS_END]
				numbers := scratchcards[i+NUMBERS_START : i+NUMBERS_END]

				winnerSet := makeWinnersSet(winners)
				myMatches := getMatches(winnerSet, numbers)
				myPoints += 1 << (myMatches - 1)
				partialMatches[j] = int(myMatches)
				j++
			}
			partialPoints <- myPoints
		}(start, end, partialMatches)
	}

	return reduceResults(partialPoints, numWorkers, matches)
}

func reduceResults(partialPoints chan int, numWorkers int, matches [][]int) common.Results[int, int] {
	var total common.Results[int, int]
	for i := 0; i < numWorkers; i++ {
		total.Part1 += <-partialPoints
	}

	for i := len(matches) - 1; i >= 0; i-- {
		for j := len(matches[i]) - 1; j >= 0; j-- {
			scratchcardsWon := 1
			myMatches := matches[i][j]
			for k, mutI, mutJ := 0, i, j; k < myMatches; k++ {
				if mutJ == len(matches[mutI])-1 {
					mutI, mutJ = mutI+1, 0
				} else {
					mutJ++
				}
				scratchcardsWon += matches[mutI][mutJ]
			}
			matches[i][j] = scratchcardsWon
			total.Part2 += scratchcardsWon
		}
	}

	return total
}

func makeWinnersSet(winners []byte) map[byte]struct{} {
	winnerSet := make(map[byte]struct{})
	var winner byte = 0
	for i := 0; i < WINNERS_LENGTH; i++ {
		n := winners[i]
		if n == ' ' {
			winnerSet[winner] = struct{}{}
			winner = 0
		} else {
			winner = winner*10 + n - '0'
		}
	}
	delete(winnerSet, 0)
	return winnerSet
}

func getMatches(set map[byte]struct{}, numbers []byte) byte {
	var number, matches byte = 0, 0
	for i := 0; i < NUMBERS_LENGTH; i++ {
		n := numbers[i]
		if n == ' ' {
			_, E := set[number]
			if E {
				matches++
			}
			number = 0
		} else {
			number = number*10 + n - '0'
		}
	}
	// Add the last number, that has '\n' after it instead of ' '
	_, E := set[number]
	if E {
		matches++
	}
	return matches
}
